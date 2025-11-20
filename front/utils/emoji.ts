// 本地表情包映射配置
export interface EmojiConfig {
  code: string; // 表情代码，如：[/微笑]
  name: string; // 表情名称，如：微笑
  filename: string; // 文件名，如：001_微笑.png
  path: string; // 表情包路径
}

// 本地表情包映射表
export const localEmojis: EmojiConfig[] = [
  { code: '[/微笑]', name: '微笑', filename: '001_微笑.png', path: '/emoji/001_微笑.png' },
  { code: '[/撇嘴]', name: '撇嘴', filename: '002_撇嘴.png', path: '/emoji/002_撇嘴.png' },
  { code: '[/色]', name: '色', filename: '003_色.png', path: '/emoji/003_色.png' },
  { code: '[/发呆]', name: '发呆', filename: '004_发呆.png', path: '/emoji/004_发呆.png' },
  { code: '[/得意]', name: '得意', filename: '005_得意.png', path: '/emoji/005_得意.png' },
  { code: '[/流泪]', name: '流泪', filename: '006_流泪.png', path: '/emoji/006_流泪.png' },
  { code: '[/害羞]', name: '害羞', filename: '007_害羞.png', path: '/emoji/007_害羞.png' },
  { code: '[/闭嘴]', name: '闭嘴', filename: '008_闭嘴.png', path: '/emoji/008_闭嘴.png' },
  { code: '[/睡]', name: '睡', filename: '009_睡.png', path: '/emoji/009_睡.png' },
  { code: '[/大哭]', name: '大哭', filename: '010_大哭.png', path: '/emoji/010_大哭.png' },
  { code: '[/尴尬]', name: '尴尬', filename: '011_尴尬.png', path: '/emoji/011_尴尬.png' },
  { code: '[/发怒]', name: '发怒', filename: '012_发怒.png', path: '/emoji/012_发怒.png' },
  { code: '[/调皮]', name: '调皮', filename: '013_调皮.png', path: '/emoji/013_调皮.png' },
  { code: '[/呲牙]', name: '呲牙', filename: '014_呲牙.png', path: '/emoji/014_呲牙.png' },
  { code: '[/惊讶]', name: '惊讶', filename: '015_惊讶.png', path: '/emoji/015_惊讶.png' },
  { code: '[/难过]', name: '难过', filename: '016_难过.png', path: '/emoji/016_难过.png' },
  { code: '[/囧]', name: '囧', filename: '017_囧.png', path: '/emoji/017_囧.png' },
  { code: '[/抓狂]', name: '抓狂', filename: '018_抓狂.png', path: '/emoji/018_抓狂.png' },
  { code: '[/吐]', name: '吐', filename: '019_吐.png', path: '/emoji/019_吐.png' },
  { code: '[/偷笑]', name: '偷笑', filename: '020_偷笑.png', path: '/emoji/020_偷笑.png' },
  { code: '[/愉快]', name: '愉快', filename: '021_愉快.png', path: '/emoji/021_愉快.png' },
  { code: '[/白眼]', name: '白眼', filename: '022_白眼.png', path: '/emoji/022_白眼.png' },
  { code: '[/傲慢]', name: '傲慢', filename: '023_傲慢.png', path: '/emoji/023_傲慢.png' },
  { code: '[/困]', name: '困', filename: '024_困.png', path: '/emoji/024_困.png' },
  { code: '[/惊恐]', name: '惊恐', filename: '025_惊恐.png', path: '/emoji/025_惊恐.png' },
  { code: '[/憨笑]', name: '憨笑', filename: '026_憨笑.png', path: '/emoji/026_憨笑.png' },
  { code: '[/悠闲]', name: '悠闲', filename: '027_悠闲.png', path: '/emoji/027_悠闲.png' },
  { code: '[/咒骂]', name: '咒骂', filename: '028_咒骂.png', path: '/emoji/028_咒骂.png' },
  { code: '[/疑问]', name: '疑问', filename: '029_疑问.png', path: '/emoji/029_疑问.png' },
  { code: '[/嘘]', name: '嘘', filename: '030_嘘.png', path: '/emoji/030_嘘.png' },
  { code: '[/晕]', name: '晕', filename: '031_晕.png', path: '/emoji/031_晕.png' },
  { code: '[/衰]', name: '衰', filename: '032_衰.png', path: '/emoji/032_衰.png' },
  { code: '[/骷髅]', name: '骷髅', filename: '033_骷髅.png', path: '/emoji/033_骷髅.png' },
  { code: '[/敲打]', name: '敲打', filename: '034_敲打.png', path: '/emoji/034_敲打.png' },
  { code: '[/再见]', name: '再见', filename: '035_再见.png', path: '/emoji/035_再见.png' },
  { code: '[/擦汗]', name: '擦汗', filename: '036_擦汗.png', path: '/emoji/036_擦汗.png' },
  { code: '[/抠鼻]', name: '抠鼻', filename: '037_抠鼻.png', path: '/emoji/037_抠鼻.png' },
  { code: '[/鼓掌]', name: '鼓掌', filename: '038_鼓掌.png', path: '/emoji/038_鼓掌.png' },
  { code: '[/坏笑]', name: '坏笑', filename: '039_坏笑.png', path: '/emoji/039_坏笑.png' },
  { code: '[/右哼哼]', name: '右哼哼', filename: '040_右哼哼.png', path: '/emoji/040_右哼哼.png' },
  { code: '[/鄙视]', name: '鄙视', filename: '041_鄙视.png', path: '/emoji/041_鄙视.png' },
  { code: '[/委屈]', name: '委屈', filename: '042_委屈.png', path: '/emoji/042_委屈.png' },
  { code: '[/快哭了]', name: '快哭了', filename: '043_快哭了.png', path: '/emoji/043_快哭了.png' },
  { code: '[/阴险]', name: '阴险', filename: '044_阴险.png', path: '/emoji/044_阴险.png' },
  { code: '[/亲亲]', name: '亲亲', filename: '045_亲亲.png', path: '/emoji/045_亲亲.png' },
  { code: '[/可怜]', name: '可怜', filename: '046_可怜.png', path: '/emoji/046_可怜.png' },
  { code: '[/笑脸]', name: '笑脸', filename: '047_笑脸.png', path: '/emoji/047_笑脸.png' },
  { code: '[/生病]', name: '生病', filename: '048_生病.png', path: '/emoji/048_生病.png' },
  { code: '[/脸红]', name: '脸红', filename: '049_脸红.png', path: '/emoji/049_脸红.png' },
  { code: '[/破涕为笑]', name: '破涕为笑', filename: '050_破涕为笑.png', path: '/emoji/050_破涕为笑.png' },
  { code: '[/恐惧]', name: '恐惧', filename: '051_恐惧.png', path: '/emoji/051_恐惧.png' },
  { code: '[/失望]', name: '失望', filename: '052_失望.png', path: '/emoji/052_失望.png' },
  { code: '[/无语]', name: '无语', filename: '053_无语.png', path: '/emoji/053_无语.png' },
  { code: '[/嘿哈]', name: '嘿哈', filename: '054_嘿哈.png', path: '/emoji/054_嘿哈.png' },
  { code: '[/捂脸]', name: '捂脸', filename: '055_捂脸.png', path: '/emoji/055_捂脸.png' },
  { code: '[/奸笑]', name: '奸笑', filename: '056_奸笑.png', path: '/emoji/056_奸笑.png' },
  { code: '[/机智]', name: '机智', filename: '057_机智.png', path: '/emoji/057_机智.png' },
  { code: '[/皱眉]', name: '皱眉', filename: '058_皱眉.png', path: '/emoji/058_皱眉.png' },
  { code: '[/耶]', name: '耶', filename: '059_耶.png', path: '/emoji/059_耶.png' },
  { code: '[/吃瓜]', name: '吃瓜', filename: '060_吃瓜.png', path: '/emoji/060_吃瓜.png' },
  { code: '[/加油]', name: '加油', filename: '061_加油.png', path: '/emoji/061_加油.png' },
  { code: '[/汗]', name: '汗', filename: '062_汗.png', path: '/emoji/062_汗.png' },
  { code: '[/天啊]', name: '天啊', filename: '063_天啊.png', path: '/emoji/063_天啊.png' },
  { code: '[/Emm]', name: 'Emm', filename: '064_Emm.png', path: '/emoji/064_Emm.png' },
  { code: '[/社会社会]', name: '社会社会', filename: '065_社会社会.png', path: '/emoji/065_社会社会.png' },
  { code: '[/旺柴]', name: '旺柴', filename: '066_旺柴.png', path: '/emoji/066_旺柴.png' },
  { code: '[/好的]', name: '好的', filename: '067_好的.png', path: '/emoji/067_好的.png' },
  { code: '[/打脸]', name: '打脸', filename: '068_打脸.png', path: '/emoji/068_打脸.png' },
  { code: '[/哇]', name: '哇', filename: '069_哇.png', path: '/emoji/069_哇.png' },
  { code: '[/翻白眼]', name: '翻白眼', filename: '070_翻白眼.png', path: '/emoji/070_翻白眼.png' },
  { code: '[/666]', name: '666', filename: '071_666.png', path: '/emoji/071_666.png' },
  { code: '[/让我看看]', name: '让我看看', filename: '072_让我看看.png', path: '/emoji/072_让我看看.png' },
  { code: '[/叹气]', name: '叹气', filename: '073_叹气.png', path: '/emoji/073_叹气.png' },
  { code: '[/苦涩]', name: '苦涩', filename: '074_苦涩.png', path: '/emoji/074_苦涩.png' },
  { code: '[/裂开]', name: '裂开', filename: '075_裂开.png', path: '/emoji/075_裂开.png' },
  { code: '[/嘴唇]', name: '嘴唇', filename: '076_嘴唇.png', path: '/emoji/076_嘴唇.png' },
  { code: '[/爱心]', name: '爱心', filename: '077_爱心.png', path: '/emoji/077_爱心.png' },
  { code: '[/心碎]', name: '心碎', filename: '078_心碎.png', path: '/emoji/078_心碎.png' },
  { code: '[/拥抱]', name: '拥抱', filename: '079_拥抱.png', path: '/emoji/079_拥抱.png' },
  { code: '[/强]', name: '强', filename: '080_强.png', path: '/emoji/080_强.png' },
  { code: '[/弱]', name: '弱', filename: '081_弱.png', path: '/emoji/081_弱.png' },
  { code: '[/握手]', name: '握手', filename: '082_握手.png', path: '/emoji/082_握手.png' },
  { code: '[/胜利]', name: '胜利', filename: '083_胜利.png', path: '/emoji/083_胜利.png' },
  { code: '[/抱拳]', name: '抱拳', filename: '084_抱拳.png', path: '/emoji/084_抱拳.png' },
  { code: '[/勾引]', name: '勾引', filename: '085_勾引.png', path: '/emoji/085_勾引.png' },
  { code: '[/拳头]', name: '拳头', filename: '086_拳头.png', path: '/emoji/086_拳头.png' },
  { code: '[/OK]', name: 'OK', filename: '087_OK.png', path: '/emoji/087_OK.png' },
  { code: '[/合十]', name: '合十', filename: '088_合十.png', path: '/emoji/088_合十.png' },
  { code: '[/啤酒]', name: '啤酒', filename: '089_啤酒.png', path: '/emoji/089_啤酒.png' },
  { code: '[/咖啡]', name: '咖啡', filename: '090_咖啡.png', path: '/emoji/090_咖啡.png' },
  { code: '[/蛋糕]', name: '蛋糕', filename: '091_蛋糕.png', path: '/emoji/091_蛋糕.png' },
  { code: '[/玫瑰]', name: '玫瑰', filename: '092_玫瑰.png', path: '/emoji/092_玫瑰.png' },
  { code: '[/凋谢]', name: '凋谢', filename: '093_凋谢.png', path: '/emoji/093_凋谢.png' },
  { code: '[/菜刀]', name: '菜刀', filename: '094_菜刀.png', path: '/emoji/094_菜刀.png' },
  { code: '[/炸弹]', name: '炸弹', filename: '095_炸弹.png', path: '/emoji/095_炸弹.png' },
  { code: '[/便便]', name: '便便', filename: '096_便便.png', path: '/emoji/096_便便.png' },
  { code: '[/月亮]', name: '月亮', filename: '097_月亮.png', path: '/emoji/097_月亮.png' },
  { code: '[/太阳]', name: '太阳', filename: '098_太阳.png', path: '/emoji/098_太阳.png' },
  { code: '[/庆祝]', name: '庆祝', filename: '099_庆祝.png', path: '/emoji/099_庆祝.png' },
  { code: '[/礼物]', name: '礼物', filename: '100_礼物.png', path: '/emoji/100_礼物.png' },
  { code: '[/红包]', name: '红包', filename: '101_红包.png', path: '/emoji/101_红包.png' },
  { code: '[/發]', name: '發', filename: '102_發.png', path: '/emoji/102_發.png' },
  { code: '[/福]', name: '福', filename: '103_福.png', path: '/emoji/103_福.png' },
  { code: '[/烟花]', name: '烟花', filename: '104_烟花.png', path: '/emoji/104_烟花.png' },
  { code: '[/爆竹]', name: '爆竹', filename: '105_爆竹.png', path: '/emoji/105_爆竹.png' },
  { code: '[/猪头]', name: '猪头', filename: '106_猪头.png', path: '/emoji/106_猪头.png' },
  { code: '[/跳跳]', name: '跳跳', filename: '107_跳跳.png', path: '/emoji/107_跳跳.png' },
  { code: '[/发抖]', name: '发抖', filename: '108_发抖.png', path: '/emoji/108_发抖.png' },
  { code: '[/转圈]', name: '转圈', filename: '109_转圈.png', path: '/emoji/109_转圈.png' },
];

// 根据表情代码获取表情配置
export function getEmojiByCode(code: string): EmojiConfig | undefined {
  return localEmojis.find(emoji => emoji.code === code);
}

// 根据表情名称获取表情配置
export function getEmojiByName(name: string): EmojiConfig | undefined {
  return localEmojis.find(emoji => emoji.name === name);
}

// 获取所有表情包配置
export function getAllEmojis(): EmojiConfig[] {
  return localEmojis;
}

// 将表情代码转换为HTML图片标签
export function convertEmojiToHtml(code: string): string {
  const emoji = getEmojiByCode(code);
  if (emoji) {
    return `<img src="${emoji.path}" alt="${emoji.name}" title="${emoji.name}" class="emoji w-5 h-5 align-middle inline mx-0.5 leading-normal">`;
  }
  return code;
}

// 解析文本中的表情代码并替换为HTML
export function parseEmojis(text: string): string {
  let result = text;
  const emojiRegex = /\[\/([^\]]+)\]/g;
  
  result = result.replace(emojiRegex, (match, emojiName) => {
    const emoji = getEmojiByName(emojiName);
    if (emoji) {
      return `<img src="${emoji.path}" alt="${emoji.name}" title="${emoji.name}" class="emoji w-5 h-5 align-middle inline mx-0.5 leading-normal">`;
    }
    return match;
  });
  
  return result;
}