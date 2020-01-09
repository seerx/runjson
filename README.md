<p align="center">
  <img width="320" src="https://github.com/seerx/runjson/blob/master/resources/logo.png">
</p>

# Run JSON
Let json object Running

# Why ?
从 Restful 转战 graphql，为了方便开发，特意封装了 graphql-go 包 <a href="https://www.github.com/seerx/goql">goql</a>，使用反射组装参数，提供额外的功能，
从开发上来说，感觉 graphql 确实比 Restful 有一些优势，比如一次请求可以调用多个 API 接口。
<br>再比如 graphql 宣传的，按需反馈、明确类型等等。
<br>刚开始接触时觉得挺不错，在使用过程中渐渐发现这些在开发过程中有些累赘，有些限制使得前后台都必须跟着调整才能适应，比如：
<br>一个 API 提供树状结构数据，使用 graphql 时，客户端必须明确指出要几层数据，这确实有些过分。
<br>其它问题不在一一罗列
<br>有鉴于此，决心做一套可以运行 json 的包，借鉴了 graphql 比较好的思想，比如一次请求多个接口。目标是让开发变得更简单。

# 目标
<ol>
    <li>让客户端开发简单易上手</li>
    <li>让服务端开发简单易上手</li>
    <li>目前使用 golang 实现一套服务端 runjson 包</li>
</ol>

# 规则
<ol>
    <li>输入 json 格式数据
        <ol>
            <li>最外层是数组</li>
            <li>数组的每一项对应一个 API 调用</li>
            <li>每一项的内容必须包含 service 项，可选包含 arg 项</li>
            <li>service 项指明 API 接口名称</li>
            <li>arg 项指明 API 接口所需参数</li>
        </ol>
    </li>
    <li>输出 json 格式数据
        <ol>
            <li>输出内容由 API 的实现决定，客户端只是获取数据</li>
            <li>输出的最外层是 dict 对象</li>
            <li>dict 对象的 key 是 API 名称</li>
            <li>dict 对象的 key 对应的 Value，是 API 执行的结果</li>
            <li>Value 是数组形式，如果同一个 API 在一起请求中，被调用 n 次，那么这个数据就有 n 维</li>
            <li>Value 数组中下标顺序，对应请求中 API 的顺序</li>
        </ol>
    </li>
    <li>提供 API 说明信息，包含 API 名称、参数说明、参数数据类型等内容</li>
</ol>

# 拓展
该包实现了 RunJSON 必要的核心功能，要用它做开发，请转到 <a href="https://www.github.com/seerx/rjhttp">rjhttp</a>，这是一个以 Run JSON 包为核心的 http 形式的服务接口包，使用它可以非常便捷的开发 API，并可以实时查看 API 文档。
