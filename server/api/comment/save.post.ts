import RPCClient from "@alicloud/pop-core";
import fs from "fs/promises";
import prisma from "~/lib/db";
import { SysConfig } from "~/lib/types";
import { sendEmail } from "~/utils/sendEmail";

type SaveCommentReq = {
  memoId: number;
  content: string;
  replyTo?: string;
  replyToId?: number;
  email?: string;
  website?: string;
  username: string;
  author: Boolean;
  token: string;
};

type recaptchaResponse = {
  success: boolean;
  challenge_ts: string;
  hostname: string;
  score: number;
  action: string;
};

const staticWord: Record<string, string> = {
  ad: "广告引流",
  political_content: "涉政内容",
  profanity: "辱骂内容",
  contraband: "违禁内容",
  sexual_content: "色情内容",
  violence: "暴恐内容",
  nonsense: "无意义内容",
  negative_content: "不良内容",
  religion: "宗教内容",
  cyberbullying: "网络暴力",
  ad_compliance: "广告法合规",
  C_customized: "违反本站规定",
};

const insertComment = async (userId: string, request: SaveCommentReq) => {
  await prisma.comment.create({
    data: {
      content: request.content,
      replyTo: request.replyTo,
      memoId: request.memoId,
      username: request.username,
      email: request.email,
      website: request.website,
      author: userId !== undefined,
    },
  });
};

const emailReg =
  /^[a-zA-Z0-9_+&*-]+(?:\.[a-zA-Z0-9_+&*-]+)*@(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,7}$/;

export default defineEventHandler(async (event) => {
  const config = ((await fs.readFile(`${process.env.CONFIG_FILE}`)).toString())
  const sysConfig = JSON.parse(config) as SysConfig
  const request = (await readBody(event)) as SaveCommentReq;
  const { content, email, memoId, replyToId, username, website, token } =
    request;

  const userId = event.context.userId;
  if (sysConfig.private.googleRecaptchaSecretKey && !token) {
    return { success: false, message: "小样儿,你是不是人机?" };
  }
  if (content.length > sysConfig.public.commentMaxLength) {
    return { success: false, message: "评论超长了,老板" };
  }
  if (username.length > 10) {
    return { success: false, message: "用户名你起这么长啥意思?" };
  }
  if (email && email.length > 30) {
    return { success: false, message: "邮箱你起这么长啥意思?" };
  }
  if (website && website.length > 100) {
    return { success: false, message: "网站地址也没这么长的啊" };
  }

  const enableAliyunTextJudge =
    sysConfig.private.enableAliyunJudge &&
    sysConfig.private.aliyunAk &&
    sysConfig.private.aliyunSk;

  if (sysConfig.private.googleRecaptchaSecretKey) {
    const url = `https://recaptcha.net/recaptcha/api/siteverify?secret=${sysConfig.private.googleRecaptchaSecretKey}&response=${token}`
    const response = (await $fetch(url)) as any as recaptchaResponse;
    if (response.score <= 0.5) {
      return {
        success: false,
        message: "小样儿,你是不是人机?",
      };
    }
  }

  if (enableAliyunTextJudge) {
    // 文本内容检查
    const aliJudgeResponse1 = (await aliTextJudge(sysConfig,
      content,
      "comment_detection"
    )) as any;
    if (
      aliJudgeResponse1.Data &&
      aliJudgeResponse1.Data.labels &&
      aliJudgeResponse1.Data.labels !== ""
    ) {
      let labelsList = aliJudgeResponse1.Data.labels.split(",");

      return {
        success: false,
        message:
          "评论内容不符合规范：" +
          labelsList.map((label: string) => staticWord[label]).join(", "),
      };
    }

    const aliJudgeResponse2 = (await aliTextJudge(sysConfig,
      username,
      "nickname_detection"
    )) as any;
    if (
      aliJudgeResponse2.Data &&
      aliJudgeResponse2.Data.labels &&
      aliJudgeResponse2.Data.labels !== ""
    ) {
      let labelsList = aliJudgeResponse2.Data.labels.split(",");

      return {
        success: false,
        message:
          "用户名不符合规范：" +
          labelsList.map((label: string) => staticWord[label]).join(", "),
      };
    }
  }

  await insertComment(userId, request);

  if (!sysConfig.private.enableNotifyByEmail) {
    // 未开启邮件通知
    return {
      success: true,
      message: "",
    };
  }
  let notificationList: any[] = [];
  notificationList.push(email || '');
  let comment = null;
  if (replyToId !== undefined && replyToId !== 0) {
    comment = await prisma.comment.findUnique({
      where: {
        id: replyToId,
      }
    });
    if (
        comment !== null &&
        comment.email !== null &&
        comment.email !== '' &&
        emailReg.test(comment.email) &&
        notificationList.indexOf(comment.email) === -1
    ){
      notificationList.push(comment.email);
      // 邮箱通知被回复者
      sendEmail({
        email: comment.email,
        subject: "新回复",
        message: `<div style="max-width: 500px;margin: 0 auto;padding: 20px;background: #f9f9f9;">
                    <div style="color: #666;line-height: 1.5;">
                      <p>您在 Moments 中的评论有回复啦！</p>
                      <p><b>${username}</b> 回复了您的评论（${comment.content}）说：</p>
                      <p style="text-indent: 30px;">${content}</p>
                      <p>请及时查阅哦~</p>
                    </div>
                    <div style="margin-top: 30px;text-align: center;font-size: 14px;color: #999;">
                      <a href="${sysConfig.public.siteUrl}/detail/${memoId}" style="margin-top: 20px;background: #007BFF;color: white;padding: 10px 20px;text-decoration: none;border-radius: 5px;">点击查看</a>
                      <p>此邮件由系统自动发送，请勿直接回复。</p>
                    </div>
                  </div>`,
      }).catch((err) => {
        console.log(`发送邮件给: ${comment.email} 失败了,原因:{$err.message}`);
      });
    }
  }

  // 非管理员
  if (event.context.userId == undefined && notificationList.indexOf(sysConfig.private.adminEmail) === -1) {
    // 判断process.env.SITE_URL是否以/结尾，如果是则去掉
    let siteUrl = sysConfig.public.siteUrl;
    if (
      siteUrl === undefined ||
      siteUrl === "" ||
      siteUrl === "undefined" ||
      siteUrl === "null"
    ) {
      siteUrl = "";
    }
    if (siteUrl.endsWith("/")) {
      siteUrl = siteUrl.slice(0, -1);
    }

    // 邮箱通知管理员
    sendEmail({
      email: sysConfig.private.adminEmail || "",
      subject: "新评论",
      message: `<div style="max-width: 500px;margin: 0 auto;padding: 20px;background: #f9f9f9;">
                  <div style="color: #666;line-height: 1.5;">
                    <p>这是一条来自 Moments 的新评论！</p>
                    <p><b>${username}</b> 说：</p>
                    <p style="text-indent: 30px;">${content}</p>
                    <p>请及时查阅哦~</p>
                  </div>
                  <div style="margin-top: 30px;text-align: center;font-size: 14px;color: #999;">
                    <a href="${siteUrl}/detail/${memoId}" style="margin-top: 20px;background: #007BFF;color: white;padding: 10px 20px;text-decoration: none;border-radius: 5px;">点击查看</a>
                    <p>此邮件由系统自动发送，请勿直接回复。</p>
                  </div>
                </div>`,
    }).catch((err) => {
      console.log(
        `发送邮箱给管理员: ${sysConfig.private.adminEmail} 失败了,原因:{$err.message}`
      );
    });
  }
  return {
    success: true,
    message: "",
  };
});

// 阿里云文本审核
async function aliTextJudge(sysConfig:SysConfig,
  content: string,
  Service: string = "comment_detection"
) {
  // 注意，此处实例化的client请尽可能重复使用，避免重复建立连接，提升检测性能。
  let client = new RPCClient({
    /**
     * 阿里云账号AccessKey拥有所有API的访问权限，建议您使用RAM用户进行API访问或日常运维。
     * 强烈建议不要把AccessKey ID和AccessKey Secret保存到工程代码里，否则可能导致AccessKey泄露，威胁您账号下所有资源的安全。
     * 常见获取环境变量方式:
     * 获取RAM用户AccessKey ID: process.env['ALIBABA_CLOUD_ACCESS_KEY_ID']
     * 获取RAM用户AccessKey Secret: process.env['ALIBABA_CLOUD_ACCESS_KEY_SECRET']
     */
    accessKeyId: sysConfig.private.aliyunAk,
    accessKeySecret: sysConfig.private.aliyunSk,
    // 接入区域和地址请根据实际情况修改
    endpoint: "https://green-cip.cn-shanghai.aliyuncs.com",
    apiVersion: "2022-03-02",
    // 设置http代理
    // httpProxy: "http://xx.xx.xx.xx:xxxx",
    // 设置https代理
    // httpsProxy: "https://username:password@xxx.xxx.xxx.xxx:9999",
  });
  // 通过以下代码创建API请求并设置参数。
  const params = {
    // 文本检测service: 内容安全控制台文本增强版规则配置的serviceCode，示例: chat_detection
    Service: Service,
    ServiceParameters: JSON.stringify({
      //待检测文本内容。
      content: content,
    }),
  };

  const requestOption = {
    method: "POST",
    formatParams: false,
  };

  let response;
  try {
    response = (await client.request(
      "TextModeration",
      params,
      requestOption
    )) as any;
    if (response.Code === 500) {
      client.endpoint = "https://green-cip.cn-beijing.aliyuncs.com";
      response = await client.request("TextModeration", params, requestOption);
    }
  } catch (err) {
    console.log(err);
    // 确保在错误情况下也返回一个值，或者抛出一个错误
    return err;
  }
  return response;
}
