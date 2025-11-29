<template>
  <Header :user="currentUser"/>

  <div class="bg-gray-100 dark:bg-gray-900 min-h-screen p-2 rounded-b-lg">


    <!-- 基本设置 -->
    <div class="bg-white dark:bg-gray-800 shadow-sm mb-2 rounded-lg overflow-hidden">
      <div class="px-4 py-2 border-b border-gray-100 dark:border-gray-700">
        <div class="flex items-center space-x-2">
          <UIcon name="i-heroicons-cog-6-tooth" class="w-4 h-4 text-gray-500 dark:text-gray-400"/>
          <h3 class="text-base font-medium text-gray-900 dark:text-white">基本</h3>
        </div>
      </div>
      <div class="divide-y divide-gray-100 dark:divide-gray-700">
        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors cursor-pointer" @click="showAdminUserName = !showAdminUserName">
          <span class="text-gray-700 dark:text-gray-300">管理账号</span>
          <div class="flex items-center space-x-2">
            <span class="text-gray-500 dark:text-gray-400 text-sm">{{ state.adminUserName }}</span>
            <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
          </div>
        </div>
        <div v-show="showAdminUserName" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600">
          <UInput v-model="state.adminUserName" size="md" class="h-10"/>
        </div>

        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors cursor-pointer" @click="showTitle = !showTitle">
          <span class="text-gray-700 dark:text-gray-300">网站标题</span>
          <div class="flex items-center space-x-2">
            <span class="text-gray-500 dark:text-gray-400 text-sm">{{ state.title }}</span>
            <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
          </div>
        </div>
        <div v-show="showTitle" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600">
          <UInput v-model="state.title" size="md" class="h-10"/>
        </div>

        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors cursor-pointer" @click="showFavicon = !showFavicon">
          <span class="text-gray-700 dark:text-gray-300">Favicon</span>
          <div class="flex items-center space-x-2">
            <UAvatar :src="state.favicon" size="sm"/>
            <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
          </div>
        </div>
        <div v-show="showFavicon" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600 space-y-3">
          <UInput v-model="state.favicon" placeholder="输入地址或上传" size="md" class="h-10"/>
          <label class="cursor-pointer inline-block">
            <UInput
              type="file"
              @change="uploadFavicon"
              accept="image/*"
              class="hidden"
            />
            <div class="flex items-center justify-center p-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition-colors text-sm">
              <UIcon name="i-heroicons-arrow-up-tray" class="w-4 h-4 mr-1"/>
              上传图标
            </div>
          </label>
        </div>
      </div>
    </div>

    <!-- 功能设置 -->
    <div class="bg-white dark:bg-gray-800 shadow-sm mb-2 rounded-lg overflow-hidden">
      <div class="px-4 py-2 border-b border-gray-100 dark:border-gray-700">
        <div class="flex items-center space-x-2">
          <UIcon name="i-heroicons-puzzle-piece" class="w-4 h-4 text-gray-500 dark:text-gray-400"/>
          <h3 class="text-base font-medium text-gray-900 dark:text-white">扩展</h3>
        </div>
      </div>
      <div class="divide-y divide-gray-100 dark:divide-gray-700">
        <div class="px-4 py-3 flex items-center justify-between">
          <div class="flex flex-col">
            <span class="text-gray-700 dark:text-gray-300">自动加载内容</span>
            <span class="text-xs text-gray-500 dark:text-gray-400">滚动时自动加载更多内容</span>
          </div>
          <UToggle v-model="state.enableAutoLoadNextPage"/>
        </div>

        <div class="px-4 py-3 flex items-center justify-between">
          <div class="flex flex-col">
            <span class="text-gray-700 dark:text-gray-300">用户主页布局</span>
            <span class="text-xs text-gray-500 dark:text-gray-400">以列表形式展示用户动态</span>
          </div>
          <UToggle v-model="state.enableNewMemo"/>
        </div>

        <div class="px-4 py-3 flex items-center justify-between">
          <div class="flex flex-col">
            <span class="text-gray-700 dark:text-gray-300">评论功能</span>
            <span class="text-xs text-gray-500 dark:text-gray-400">允许用户对动态进行评论</span>
          </div>
          <UToggle v-model="state.enableComment"/>
        </div>

        <div class="px-4 py-3 flex items-center justify-between">
          <div class="flex flex-col">
            <span class="text-gray-700 dark:text-gray-300">开放注册</span>
            <span class="text-xs text-gray-500 dark:text-gray-400">允许新用户注册账号</span>
          </div>
          <UToggle v-model="state.enableRegister"/>
        </div>
      </div>
    </div>

    <!-- 显示设置 -->
    <div class="bg-white dark:bg-gray-800 shadow-sm mb-2 rounded-lg overflow-hidden">
      <div class="px-4 py-2 border-b border-gray-100 dark:border-gray-700">
        <div class="flex items-center space-x-2">
          <UIcon name="i-heroicons-eye" class="w-4 h-4 text-gray-500 dark:text-gray-400"/>
          <h3 class="text-base font-medium text-gray-900 dark:text-white">显示</h3>
        </div>
      </div>
      <div class="divide-y divide-gray-100 dark:divide-gray-700">
        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors cursor-pointer" @click="showBeiAnNo = !showBeiAnNo">
          <span class="text-gray-700 dark:text-gray-300">备案号</span>
          <div class="flex items-center space-x-2">
            <span class="text-gray-500 dark:text-gray-400 text-sm">{{ state.beiAnNo || "未设置" }}</span>
            <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
          </div>
        </div>
        <div v-show="showBeiAnNo" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600">
          <UInput v-model="state.beiAnNo" placeholder="没有可以不填写" size="md" class="h-10"/>
        </div>

        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors cursor-pointer" @click="showRssMaxItems = !showRssMaxItems">
          <span class="text-gray-700 dark:text-gray-300">RSS最大条数</span>
          <div class="flex items-center space-x-2">
            <span class="text-gray-500 dark:text-gray-400 text-sm">{{ state.rssMaxItems || "15(默认)" }}</span>
            <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
          </div>
        </div>
        <div v-show="showRssMaxItems" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600">
          <UInput v-model="rssMaxItemsComputed" placeholder="留空则默认(15条)" size="md" class="h-10"/>
        </div>

        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors cursor-pointer" @click="showMaxCommentLength = !showMaxCommentLength">
          <span class="text-gray-700 dark:text-gray-300">评论最大字数</span>
          <div class="flex items-center space-x-2">
            <span class="text-gray-500 dark:text-gray-400 text-sm">{{ state.maxCommentLength }}</span>
            <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
          </div>
        </div>
        <div v-show="showMaxCommentLength" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600">
          <UInput v-model.number="state.maxCommentLength" size="md" class="h-10"/>
        </div>

        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors cursor-pointer" @click="showMemoMaxHeight = !showMemoMaxHeight">
          <span class="text-gray-700 dark:text-gray-300">发言最大高度</span>
          <div class="flex items-center space-x-2">
            <span class="text-gray-500 dark:text-gray-400 text-sm">{{ state.memoMaxHeight ? state.memoMaxHeight + 'px' : '不限制' }}</span>
            <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
          </div>
        </div>
        <div v-show="showMemoMaxHeight" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600">
          <UInput v-model.number="state.memoMaxHeight" placeholder="单位px,填0时则不限制高度" size="md" class="h-10"/>
        </div>

        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors cursor-pointer" @click="showCommentOrder = !showCommentOrder">
          <span class="text-gray-700 dark:text-gray-300">评论排序方式</span>
          <div class="flex items-center space-x-2">
            <span class="text-gray-500 dark:text-gray-400 text-sm">{{ state.commentOrder === 'desc' ? '倒序' : '正序' }}</span>
            <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
          </div>
        </div>
        <div v-show="showCommentOrder" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600">
          <USelectMenu v-model="state.commentOrder"
                       :options="[{label:'倒序,越晚发布越靠前',value:'desc'},{label:'正序,越早发布越靠前',value:'asc'}]"
                       value-attribute="value" option-attribute="label" size="md" class="w-full"/>
        </div>

        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors cursor-pointer" @click="showTimeFormat = !showTimeFormat">
          <span class="text-gray-700 dark:text-gray-300">日期格式</span>
          <div class="flex items-center space-x-2">
            <span class="text-gray-500 dark:text-gray-400 text-sm">{{ state.timeFormat === 'timeAgo' ? '几分钟前' : $dayjs().format('YYYY-MM-DD HH:mm') }}</span>
            <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
          </div>
        </div>
        <div v-show="showTimeFormat" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600">
          <USelectMenu v-model="state.timeFormat"
                       :options="[{label:'几分钟前',value:'timeAgo'},{label:$dayjs().format('YYYY-MM-DD HH:mm'),value:'time'}]"
                       value-attribute="value" option-attribute="label" size="md" class="w-full"/>
        </div>
      </div>
    </div>

    <!-- 安全设置 -->
    <div class="bg-white dark:bg-gray-800 shadow-sm mb-2 rounded-lg overflow-hidden">
      <div class="px-4 py-2 border-b border-gray-100 dark:border-gray-700">
        <div class="flex items-center space-x-2">
          <UIcon name="i-heroicons-shield-check" class="w-4 h-4 text-gray-500 dark:text-gray-400"/>
          <h3 class="text-base font-medium text-gray-900 dark:text-white">安全</h3>
        </div>
      </div>
      <div class="divide-y divide-gray-100 dark:divide-gray-700">
        <div class="px-4 py-3 flex items-center justify-between">
          <span class="text-gray-700 dark:text-gray-300">启用Google Recaptcha</span>
          <UToggle v-model="state.enableGoogleRecaptcha"/>
        </div>

        <template v-if="state.enableGoogleRecaptcha">
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="SiteKey" name="googleSiteKey" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.googleSiteKey" size="md" class="h-10"/>
            </UFormGroup>
          </div>
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="SecretKey" name="googleSecretKey" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.googleSecretKey" type="password" size="md" class="h-10"/>
            </UFormGroup>
          </div>
        </template>
      </div>
    </div>

    <!-- OIDC设置 -->
    <div class="bg-white dark:bg-gray-800 shadow-sm mb-2 rounded-lg overflow-hidden">
      <div class="px-4 py-2 border-b border-gray-100 dark:border-gray-700">
        <div class="flex items-center space-x-2">
          <UIcon name="i-heroicons-key" class="w-4 h-4 text-gray-500 dark:text-gray-400"/>
          <h3 class="text-base font-medium text-gray-900 dark:text-white">OIDC认证</h3>
        </div>
      </div>
      <div class="divide-y divide-gray-100 dark:divide-gray-700">
        <div class="px-4 py-3 flex items-center justify-between">
          <div class="flex flex-col">
            <span class="text-gray-700 dark:text-gray-300">启用OIDC认证</span>
            <span class="text-xs text-gray-500 dark:text-gray-400">使用OIDC标准协议进行用户认证</span>
          </div>
          <UToggle v-model="state.enableOIDC"/>
        </div>

        <template v-if="state.enableOIDC">
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="Issuer" name="oidcIssuer" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.oidcIssuer" placeholder="https://accounts.example.com" size="md" class="h-10"/>
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">OIDC提供商的Issuer URL</p>
            </UFormGroup>
          </div>
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="Client ID" name="oidcClientId" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.oidcClientId" size="md" class="h-10"/>
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">在OIDC提供商处注册的客户端ID</p>
            </UFormGroup>
          </div>
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="Client Secret" name="oidcClientSecret" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.oidcClientSecret" type="password" size="md" class="h-10"/>
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">在OIDC提供商处注册的客户端密钥</p>
            </UFormGroup>
          </div>
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="Redirect URI" name="oidcRedirectUri" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.oidcRedirectUri" placeholder="https://your-site.com/api/auth/oidc/callback" size="md" class="h-10"/>
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">OIDC认证后的回调地址</p>
            </UFormGroup>
          </div>
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="Scopes" name="oidcScopes" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.oidcScopes" placeholder="openid profile email" size="md" class="h-10"/>
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">请求的OIDC作用域，空格分隔</p>
            </UFormGroup>
          </div>
        </template>
      </div>
    </div>

    <!-- 存储设置 -->
    <div class="bg-white dark:bg-gray-800 shadow-sm mb-2 rounded-lg overflow-hidden">
      <div class="px-4 py-2 border-b border-gray-100 dark:border-gray-700">
        <div class="flex items-center space-x-2">
          <UIcon name="i-heroicons-server" class="w-4 h-4 text-gray-500 dark:text-gray-400"/>
          <h3 class="text-base font-medium text-gray-900 dark:text-white">存储</h3>
        </div>
      </div>
      <div class="divide-y divide-gray-100 dark:divide-gray-700">
        <div class="px-4 py-3 flex items-center justify-between">
          <span class="text-gray-700 dark:text-gray-300">启用S3存储</span>
          <UToggle v-model="state.enableS3"/>
        </div>

        <template v-if="state.enableS3">
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="Bucket 域名（资源访问地址）" name="domain" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.s3.domain" placeholder="https://moments-test-bucket.oss-cn-hangzhou.aliyuncs.com" size="md" class="h-10"/>
            </UFormGroup>
          </div>
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="Endpoint 地址" name="endpoint" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.s3.endpoint" placeholder="https://oss-cn-hangzhou.aliyuncs.com" size="md" class="h-10"/>
            </UFormGroup>
          </div>
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="Bucket 名称" name="bucket" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.s3.bucket" placeholder="moments-test-bucket" size="md" class="h-10"/>
            </UFormGroup>
          </div>
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="Bucket 地区" name="region" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.s3.region" placeholder="oss-cn-hangzhou" size="md" class="h-10"/>
            </UFormGroup>
          </div>
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="AccessKey" name="accessKey" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.s3.accessKey" size="md" class="h-10"/>
            </UFormGroup>
          </div>
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="SecretKey" name="secretKey" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.s3.secretKey" type="password" size="md" class="h-10"/>
            </UFormGroup>
          </div>
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="图片后缀（在访问缩略图时会追加在图片地址后）" name="thumbnailSuffix" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.s3.thumbnailSuffix" size="md" class="h-10"/>
            </UFormGroup>
          </div>
        </template>
      </div>
    </div>

    <!-- 邮件设置 -->
    <div class="bg-white dark:bg-gray-800 shadow-sm mb-2 rounded-lg overflow-hidden">
      <div class="px-4 py-2 border-b border-gray-100 dark:border-gray-700">
        <div class="flex items-center space-x-2">
          <UIcon name="i-heroicons-envelope" class="w-4 h-4 text-gray-500 dark:text-gray-400"/>
          <h3 class="text-base font-medium text-gray-900 dark:text-white">邮件</h3>
        </div>
      </div>
      <div class="divide-y divide-gray-100 dark:divide-gray-700">
        <div class="px-4 py-3 flex items-center justify-between">
          <span class="text-gray-700 dark:text-gray-300">启用邮件通知</span>
          <UToggle v-model="state.enableEmail"/>
        </div>

        <template v-if="state.enableEmail">
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="smtp服务器" name="smtpHost" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.smtpHost" placeholder="smtp.qq.com" size="md" class="h-10"/>
            </UFormGroup>
          </div>
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="smtp端口" name="smtpPort" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.smtpPort" placeholder="465" size="md" class="h-10"/>
            </UFormGroup>
          </div>
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="smtp用户名" name="smtpUsername" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.smtpUsername" placeholder="******@qq.com" size="md" class="h-10"/>
            </UFormGroup>
          </div>
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="smtp密码/授权码" name="smtpPassword" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.smtpPassword" type="password" size="md" class="h-10"/>
            </UFormGroup>
          </div>
        </template>
      </div>
    </div>

    <!-- 企业微信设置 -->
    <div class="bg-white dark:bg-gray-800 shadow-sm mb-2 rounded-lg overflow-hidden">
      <div class="px-4 py-2 border-b border-gray-100 dark:border-gray-700">
        <div class="flex items-center space-x-2">
          <UIcon name="i-simple-icons-wechat" class="w-4 h-4 text-gray-500 dark:text-gray-400"/>
          <h3 class="text-base font-medium text-gray-900 dark:text-white">企业微信</h3>
        </div>
      </div>
      <div class="divide-y divide-gray-100 dark:divide-gray-700">
        <div class="px-4 py-3 flex items-center justify-between">
          <div class="flex flex-col">
            <span class="text-gray-700 dark:text-gray-300">企业微信通知</span>
            <span class="text-xs text-gray-500 dark:text-gray-400">有评论/回复时，能及时推送通知</span>
          </div>
          <UToggle v-model="state.enableWechatWebhook"/>
        </div>

        <template v-if="state.enableWechatWebhook">
          <div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <UFormGroup label="Webhook Key" name="wechatWebhookUrl" :ui="{label:{base:'font-bold'}}" class="w-full">
              <UInput v-model="state.wechatWebhookUrl" placeholder="73a***-***-***-***-***" size="md" class="h-10"/>
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">填写企业微信群机器人 Webhook URL 中 key = 后的值</p>
            </UFormGroup>
          </div>
        </template>
      </div>
    </div>

    <!-- 高级设置 -->
    <div class="bg-white dark:bg-gray-800 shadow-sm mb-2 rounded-lg overflow-hidden">
      <div class="px-4 py-2 border-b border-gray-100 dark:border-gray-700">
        <div class="flex items-center space-x-2">
          <UIcon name="i-heroicons-wrench-screwdriver" class="w-4 h-4 text-gray-500 dark:text-gray-400"/>
          <h3 class="text-base font-medium text-gray-900 dark:text-white">高级</h3>
        </div>
      </div>
      <div class="divide-y divide-gray-100 dark:divide-gray-700">
        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors cursor-pointer" @click="showCss = !showCss">
          <span class="text-gray-700 dark:text-gray-300">自定义CSS</span>
          <div class="flex items-center space-x-2">
            <span class="text-gray-500 dark:text-gray-400 text-sm">{{ state.css ? "已设置" : "未设置" }}</span>
            <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
          </div>
        </div>
        <div v-show="showCss" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600">
          <UTextarea v-model="state.css" :rows="5" placeholder="输入自定义CSS代码"/>
        </div>

        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors cursor-pointer" @click="showJs = !showJs">
          <span class="text-gray-700 dark:text-gray-300">自定义JS</span>
          <div class="flex items-center space-x-2">
            <span class="text-gray-500 dark:text-gray-400 text-sm">{{ state.js ? "已设置" : "未设置" }}</span>
            <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
          </div>
        </div>
        <div v-show="showJs" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600">
          <UTextarea v-model="state.js" :rows="5" placeholder="输入自定义JS代码"/>
        </div>
      </div>
    </div>

    <!-- 版本信息 -->
    <div class="bg-white dark:bg-gray-800 shadow-sm mb-2 rounded-lg overflow-hidden">
      <div class="px-4 py-2 border-b border-gray-100 dark:border-gray-700">
        <div class="flex items-center space-x-2">
          <UIcon name="i-weui-info-outlined" class="w-4 h-4 text-gray-500 dark:text-gray-400"/>
          <h3 class="text-base font-medium text-gray-900 dark:text-white">关于</h3>
        </div>
      </div>
      <div class="px-4 py-3 space-y-3">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 bg-gradient-to-br from-blue-400 to-blue-600 rounded-lg flex items-center justify-center">
            <span class="text-white font-bold text-lg">M</span>
          </div>
          <div class="flex flex-col">
            <span class="text-gray-700 dark:text-gray-300 font-medium">极简朋友圈</span>
            <span class="text-xs text-gray-500 dark:text-gray-400">一个极简、开源的朋友圈应用</span>
          </div>
        </div>
        <div class="space-y-2 pt-2 border-t border-gray-100 dark:border-gray-700">
          <div v-if="version" class="flex items-center justify-between">
            <span class="text-sm text-gray-600 dark:text-gray-300">版本号</span>
            <span class="text-sm text-gray-500 dark:text-gray-400">{{ version }}</span>
          </div>
          <div v-if="commitId" class="flex items-center justify-between">
            <span class="text-sm text-gray-600 dark:text-gray-300">CommitID</span>
            <span class="text-sm text-gray-500 dark:text-gray-400 font-mono">{{ commitId }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 保存按钮 -->
    <div class="py-3">
      <UButton class="w-full justify-center bg-blue-500 hover:bg-blue-600" @click="save" size="md">保存设置</UButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import type {SysConfigVO, UserVO} from "~/types";
import {toast} from "vue-sonner";
import {useUpload} from "~/utils";

const currentUser = useState<UserVO>('userinfo')
const version = ref('')
const commitId = ref('')

// 控制各个设置项的显示/隐藏
const showAdminUserName = ref(false)
const showTitle = ref(false)
const showFavicon = ref(false)
const showBeiAnNo = ref(false)
const showRssMaxItems = ref(false)
const showMaxCommentLength = ref(false)
const showMemoMaxHeight = ref(false)
const showCommentOrder = ref(false)
const showTimeFormat = ref(false)
const showCss = ref(false)
const showJs = ref(false)

const state = reactive({
  enableGoogleRecaptcha: false,
  googleSiteKey:"",
  googleSecretKey:"",
  enableAutoLoadNextPage: true,
  enableNewMemo: true,
  enableComment: true,
  enableRegister: true,
  maxCommentLength: 120,
  memoMaxHeight: 300,
  commentOrder: 'desc',
  timeFormat: 'timeAgo',
  adminUserName: "admin",
  title: "极简朋友圈",
  favicon: "/favicon.ico",
  beiAnNo: "",
  css: "",
  js: "",
  rssMaxItems: 15,
  enableS3: false,
  s3: {
    domain: "",
    bucket: "",
    region: "",
    accessKey: "",
    secretKey: "",
    endpoint: "",
    thumbnailSuffix: ""
  },
  enableEmail: false,
  smtpHost: "",
  smtpPort: "",
  smtpUsername: "",
  smtpPassword: "",
  enableWechatWebhook: false,
  wechatWebhookUrl: "",
  enableOIDC: false,
  oidcIssuer: "",
  oidcClientId: "",
  oidcClientSecret: "",
  oidcRedirectUri: "",
  oidcScopes: "openid profile email",
})

// 处理 rssMaxItems，确保其始终为数字
const rssMaxItemsComputed = computed({
  get: () => state.rssMaxItems,
  set: (val: number | string | null | undefined) => {
    if (val === null || val === undefined || val === '') {
      state.rssMaxItems = 0;
    } else {
      state.rssMaxItems = Number(val);
    }
  }
});

const reload = async () => {
  const res = await useMyFetch<SysConfigVO>('/sysConfig/getFull')
  if (res) {
    Object.assign(state, res)
    version.value = res.version
    commitId.value = res.commitId
  }
}

const save = async () => {
  await useMyFetch('/sysConfig/save', state)
  toast.success("保存成功")
  location.reload()
}

const uploadFavicon = async (files: FileList) => {
  for (let i = 0; i < files.length; i++) {
    if (files[i].type.indexOf("image") < 0){
      toast.error("只能上传图片");
      return
    }
  }
  const result = await useUpload(files)
  if (result.length) {
    toast.success("上传成功")
    state.favicon = result[0]
  }
}

onMounted(async () => {
  await reload()
})

</script>
