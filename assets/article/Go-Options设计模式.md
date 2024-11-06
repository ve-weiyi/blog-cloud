作为 Golang 开发人员，遇到的众多问题之一是试图将函数的参数设为可选。这是一个非常常见的用例，有一些对象应该使用一些基本的默认设置开箱即用，并且您可能偶尔想要提供一些更详细的配置。
在 python、kotlin 中，你可以给参数一个默认值，并在调用方法时省略它们。但是在 Golang 中，是无法这么做。
那么怎么解决这个问题呢？ 答案就是Options模式。Options模式在Golang中应用十分广泛，几乎每个框架中都有他的踪迹。
引言
在 Go 语言中，函数选项模式（Function Options Pattern）是一种常见且强大的设计模式，用于构建可扩展、易于使用和灵活的 API。该模式允许开发人员通过函数参数选项的方式来配置和定制函数的行为，从而避免函数参数过多和复杂的问题。本文将从多个方面详细介绍函数选项模式的实现原理、使用场景和具体实例，帮助大家全面理解和应用这一设计模式。
设计目标

假设让你设计一个功能，可以实现http调用，并且可以设置超时时间、请求header、请求param、请求body等参数，你会怎么做？
你可能会这样。这是一个简单的函数传参
func HttpRequestGET(url string, headers map[string]string, params map[string]string, body []byte) ([]byte, error) {

    return content, nil
}

func HttpRequestPOST(url string, headers map[string]string, params map[string]string, body []byte) ([]byte, error) {

    return content, nil
}
- 参数都是必传的。
- 参数的顺序是固定的。
- 新增一个参数，需要修改所有的方法。
- 参数很多时，这个函数的入参列表会变得很多。
- 相同类型的方法，需要传入重复地传入相同参数。

如何优化？
1. 将所有的参数存储在一个结构体中，在需要的时候进行使用。
2. 给每一个参数配置默认的值，不需要的参数不用传入。
// ClientBuilder 用于构建 HTTP 客户端。
type HttpClient struct {
    timeout    time.Duration // 请求超时时间

    method  string
    headers map[string]string // 请求头部
    params  map[string]string // 请求参数
    body    []byte            // 请求体
}

func NewHttpClient(timeout time.Duration, method string, headers, params map[string]string, body []byte) *HttpClient {
        h := HttpClient{}
        h.timeout = timeout
        h.method = method
        h.headers = headers
        h.params = params
        h.body = body
        return &h
}

设计实现
如果你熟悉Java，你会知道Java有23种设计模式，如单例模式、工厂模式、建造者模式、观察者模式...
对于以上需求，你可能会想到使用建造者模式实现。
builder
Builder 模式旨在通过链式调用或者顺序调用一系列方法来构建复杂对象。它将对象的构建与其表示分离，允许用户以可读性更高的方式构建对象。
优点：
1. 可读性：通过链式调用或者顺序调用方法，使得代码易于理解和维护。
2. 类型安全：Builder 模式可以在构建过程中强制执行类型检查，减少错误。
3. 隔离构建逻辑：将构建逻辑与对象表示分离，使得构建过程更加灵活且可扩展。
缺点：
1. 实现复杂度：Builder 模式的实现可能会增加代码复杂度，尤其是当需要构建的对象结构较为复杂时。
2. 不够灵活：相比于函数选项模式，Builder 模式可能不够灵活，因为用户需要按照预定义的步骤构建对象。

需求
golang  使用Builder模式设计封装http库，使这个方法可以设置超时时间、请求header、请求param、请求body。
代码
package httpx

import (
    "bytes"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
    "time"
)

// ClientBuilder 用于构建 HTTP 客户端。
type ClientBuilder struct {
    httpClient *http.Client  // 底层HTTP客户端
    timeout    time.Duration // 请求超时时间

    headers map[string]string // 请求头部
    params  map[string]string // 请求参数
    body    []byte            // 请求体
}

// NewClientBuilder 创建一个具有默认设置的新 ClientBuilder。
func NewClientBuilder() *ClientBuilder {
    return &ClientBuilder{
       httpClient: &http.Client{},
       timeout:    30 * time.Second, // 默认超时时间
       headers:    make(map[string]string),
       params:     make(map[string]string),
    }
}

// WithTimeout 设置 HTTP 请求的超时时间。
func (c *ClientBuilder) WithTimeout(timeout time.Duration) *ClientBuilder {
    c.timeout = timeout
    return c
}

// WithHeaders 设置 HTTP 请求的头部。
func (c *ClientBuilder) WithHeaders(headers map[string]string) *ClientBuilder {
    c.headers = headers
    return c
}

// WithParams 设置 HTTP 请求的参数。
func (c *ClientBuilder) WithParams(params map[string]string) *ClientBuilder {
    c.params = params
    return c
}

// WithBodyObject 设置 HTTP 请求的正文。
func (c *ClientBuilder) WithBody(body []byte) *ClientBuilder {
    c.body = body
    return c
}

// Build 根据构建器的设置创建一个新的 HTTP 客户端。
func (c *ClientBuilder) Build() *http.Client {
    client := &http.Client{
       Timeout: c.timeout,
    }

    return client
}

// DoRequest 执行一个 HTTP 请求。
func (c *ClientBuilder) DoRequest(method, rawURL string) (respBody []byte, err error) {
    uv, err := url.ParseRequestURI(rawURL)
    if err != nil {
       return nil, err
    }

    // 创建请求
    req, err := http.NewRequest(method, uv.String(), nil)
    if err != nil {
       return nil, err
    }

    // 设置头部
    for key, value := range c.headers {
       req.Header.Set(key, value)
    }

    // 设置查询参数
    query := req.URL.Query()
    for key, value := range c.params {
       query.Add(key, value)
    }
    req.URL.RawQuery = query.Encode()

    // 设置请求体
    req.Body = io.NopCloser(bytes.NewReader(c.body))

    log.Println("requestUrl:", req.URL.String())
    // 执行请求
    resp, err := c.httpClient.Do(req)
    if err != nil {
       return nil, err
    }

    // 检查响应状态码
    if resp.StatusCode != http.StatusOK {
       return nil, fmt.Errorf("http request fail. url:%v, code:%d,err:%s", req.URL.String(), resp.StatusCode, resp.Status)
    }

    defer resp.Body.Close()
    respBody, err = io.ReadAll(resp.Body)
    if err != nil {
       return nil, err
    }

    return respBody, nil
}
options
函数选项模式基于函数参数的可变性和可选性来实现。它通过将函数的配置选项作为参数传递给函数，从而实现了函数行为的定制。通过使用函数选项模式，我们可以避免创建大量的函数重载或参数组合，提高代码的可读性和可维护性。
函数选项模式的实现依赖于 Go 语言的可变参数和函数类型。在 Go 语言中，我们可以使用可变参数来接收不定数量的函数选项，并将这些选项保存在一个结构体中。结构体的字段可以存储选项的值，而字段的类型可以是函数类型，用于执行选项所需的操作。通过将选项存储在结构体中，我们可以在函数内部轻松地访问和使用这些选项。
优点：
1. 简单直观：函数选项模式易于理解和实现。
2. 灵活性：用户可以根据需要选择性地传递参数，更改函数或结构体的行为。
3. 易于扩展：添加新选项时，不需要修改现有的代码，只需要定义新的函数选项。
缺点：
1. 可读性：在一些情况下，函数调用可能变得复杂，难以阅读和理解。
需求
golang  使用Option模式设计封装http库，使这个方法可以设置超时时间、请求header、请求param、请求body。
代码
package httpx

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    "time"
)

// Client 表示HTTP客户端。
type Client struct {
    httpClient *http.Client  // 底层HTTP客户端
    timeout    time.Duration // 请求超时时间

    headers map[string]string // 请求头部
    params  map[string]string // 请求参数
    body    []byte            // 请求体
}

// Option 表示用于配置HTTP客户端的函数选项。
type Option func(*Client)

// NewClient 使用默认设置创建一个新的HTTP客户端。
func NewClient(options ...Option) *Client {
    client := &Client{
       httpClient: &http.Client{},
       timeout:    30 * time.Second, // 默认超时时间
       headers:    make(map[string]string),
       params:     make(map[string]string),
    }

    // 应用选项
    for _, option := range options {
       option(client)
    }

    return client
}

// WithTimeout 设置HTTP请求的超时时间。
func WithTimeout(timeout time.Duration) Option {
    return func(c *Client) {
       c.timeout = timeout
    }
}

// WithHeaders 设置HTTP请求的头部。
func WithHeaders(headers map[string]string) Option {
    return func(c *Client) {
       for key, value := range headers {
          c.headers[key] = value
       }
    }
}

// WithParams 设置HTTP请求的参数。
func WithParams(params map[string]string) Option {
    return func(c *Client) {
       for key, value := range params {
          c.params[key] = value
       }
    }
}

// WithBody 设置HTTP请求的请求体。
func WithBody(body []byte) Option {
    return func(c *Client) {
       c.body = body
    }
}

// DoRequest 执行一个HTTP请求。
func (c *Client) DoRequest(method, url string) (respBody []byte, err error) {
    // 使用请求体创建请求
    req, err := http.NewRequest(method, url, bytes.NewBuffer(c.body))
    if err != nil {
       return nil, err
    }

    // 设置头部
    for key, value := range c.headers {
       req.Header.Set(key, value)
    }

    // 设置查询参数
    query := req.URL.Query()
    for key, value := range c.params {
       query.Add(key, value)
    }
    req.URL.RawQuery = query.Encode()

    // 设置超时时间
    c.httpClient.Timeout = c.timeout

    // 执行请求
    resp, err := c.httpClient.Do(req)
    if err != nil {
       return nil, err
    }

    // 检查响应状态码
    if resp.StatusCode != http.StatusOK {
       return nil, fmt.Errorf("http request fail. url:%v, code:%d,err:%s", req.URL.String(), resp.StatusCode, resp.Status)
    }

    defer resp.Body.Close()
    respBody, err = io.ReadAll(resp.Body)
    if err != nil {
       return nil, err
    }

    return respBody, nil
}

使用方式
package httpx

import (
    "testing"
    "time"
)

func TestNewClientBuilder(t *testing.T) {
    // Create a new HTTP httpClient with options using Builder pattern
    resp, err := NewClientBuilder().
       WithTimeout(10*time.Second).
       WithHeaders(map[string]string{"Content-Type": "application/json"}).
       WithParams(map[string]string{"param1": "value1", "param2": "value2"}).
       WithBody([]byte(`{"key": "value"}`)).
       DoRequest("GET", "https://baidu.com")

    t.Log(string(resp), err)
}

func TestNewClientOptions(t *testing.T) {
    // Create a new HTTP httpClient with options using Option pattern
    resp, err := NewClient(
       WithTimeout(10*time.Second),
       WithHeaders(map[string]string{"Content-Type": "application/json"}),
       WithParams(map[string]string{"param1": "value1", "param2": "value2"}),
       WithBody([]byte(`{"key": "value"}`)),
    ).DoRequest("GET", "https://baidu.com")

    t.Log(string(resp), err)
}
结论
推荐
- Options 模式在封装库很常被使用，将一些功能封装成对象，使其支持多个可选参数。Options 模式比 Builder 模式简洁且对于参数比较少的对象使用更方便。但是对于有许多参数的对象就会很啰嗦
- Builder 模式允许创建具有许多可选参数的复杂对象。它将对象的构造与其表示分开，并提供了一种使用相同构造过程创建同一对象的不同表示的方法。相比较建造者模式会更加强大。
问题
- Golang 还有哪些设计模式？
