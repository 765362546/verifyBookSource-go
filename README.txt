
参考 https://github.com/WuSuoV/verifyBookSource 写一个go语言版本的书源校验工具，作为学习go语言的练手项目

1. 判断类型 是url还是本地文件
2. 解析成对象
3. 处理逻辑
   获取bookSourceUrl,去重，get，超时时间，成功，加上true，失败，加上false
4. 并发处理
   将书源解析成切片，投到channel中


分解:

原始书源列表--切片
    书源  map对象


0. 程序入口
   判断从config.json读取，还是从命令行读取
   config.json ---- 解析为config对象
   命令行       ---- 解析为config对象 

   config对象:
      path:  书源地址，可能为url，也可能为file路径
      workers: 线程数 即channel大小
      timeout: 超时时间 默认3s
      outpath： 输出路径，默认当前


1. 判断类型
   url  -- 下载 ---解析为书源列表
   file-----------解析为书源列表

2. 生产者A，遍历书源列表，投入到channel中
去重逻辑 -- TODO

3. 消费者B，从channel中读取单个书源，get + 超时, 成功的，投入到成功的channel，失败的，投入到失败的channel --- 直接追加到各自的切片中
4. 验证后的书源  拼接成新的 书源列表----写入到文件中


5. 统计结果
书源总数
有效书源数 计算得出
无效书源数 计算得出
重复书源数

耗时  计算得出


----------------


