## 关于博客rpc的设计

https://blog.51cto.com/u_14201949/6792978

在最初的springboot版本里，整个博客服务是一个单体服务，所有数据都使用一张表存储。只需要运行一个java服务，就可以完成整个博客的功能。

在gin版本里，沿用之前的设计，整个服务仍然是一个单体服务，只需要运行一个go服务，所有数据都使用一张表存储。就可以完成整个博客的功能。

当服务体量变大时，单体服务的缺点就暴露出来了，单体服务的缺点有：
一个服务挂，整个挂。

go-zero是一个微服务框架，它的设计思想是将一个大的服务拆分成多个小的服务。
每一个小的服务都可以是一个独立的服务，每一个服务都有自己的数据库，自己的rpc服务，自己的api服务。

每个小服务都可以拥有自己的数据库，如 UserRpc 只需要用户相关的表，并且 UserRpc 可以提供给 ArticleRpc 使用，也可以提供给 Comment 使用。
还可以提供给其他不是博客服务的服务使用。服务可以更好的解耦。

Q：我有两个应用，一个是blog博客服务，另一个是shop商店服务，我希望两个应用都使用同一个用户数据库，我应该怎么设计服务？
A：当然，你可以把所有的表都堆到一个数据库，这样这个数据库会非常大，并且很杂，很显然这样做很low，这样的设计是不合理的。
你可以设计三个数据库：user、blog、shop，user数据库只存储用户相关的表，blog数据库只存储博客相关的表，shop数据库只存储商店相关的表。
启动三个rpc服务，userRpc、blogRpc、shopRpc，blog和shop服务都可以调用user服务。
启动两个api服务，blogApi、shopApi，blogApi只提供给blog服务调用，shopApi只提供给shop服务调用。

Q：既然只需要设计两个Rpc服务


### v1版本

最初的版本，我把每一个数据库操作当中一个rpc服务，这样的好处是可以很好的解耦，但是这样的缺点是每一个rpc服务都要写一个proto文件，这样的工作量太大了。
每一个rpc只是简单的增删改查操作，不能与其他rpc服务进行交互。
实际上rpc的功能与repository层是垂直的、重合的，这样的设计是不合理的。

### v2版本
我把rpc分成了几个模块，各模块直接是平级关系，每一个模块都可以单独完成一个小功能的运转。
例如article模块可以完成文章的增删改查，并且可以完成文章与标签、分类的关联。


当前的服务结构如下：

- account 
  - user
  - user_login_history

- permission
  - role
  - menu
  - api
  
- article
  - article
  - article_category
  - article_tag
  
- message 
 - comment
 - remark
 - chat

- photo
  - album
  - photo

- resource
  - file 

- website
  - visit
  - page
  - talk
  - friend
  
- config
  - config

- syslog
  - operation_log
  - upload_log

什么事微服务，微服务是一种架构风格，它是一种将一个应用程序设计为一组小的服务的方法，每个服务都可以独立部署、独立运行、独立维护。
例如，一个博客服务，可以拆分成用户服务、文章服务、评论服务、权限服务、上传服务、网站服务、日志服务等等。
这些功能对应一个个rpc服务，每个rpc服务都可以独立运行、独立部署、独立维护。
如何设计rpc。我认为，每个rpc应该自成一个模块，可以独立完成工作。
例如，article rpc可以完成文章的增删改查，可以完成文章与标签、分类的关联。
那为什么不把article rpc的功能合并到user rpc里呢？因为这样的设计是不合理的，user rpc应该只完成用户相关的工作，article rpc应该只完成文章相关的工作。
那为什么不把article rpc的功能拆分成article rpc、tag rpc、category rpc呢？因为article和tag、category共同完成一个功能，拆分后article rpc需要调用tag rpc、category rpc，才能实现功能的完整。
那为什么不把article rpc、user rpc都合并成一个blog rpc呢，这样共同组成一个博客服务。当然，这是可以的。如何拆分，拆分到多细，取决于个人的想法。
作为设计者，你应该思考，如果把这个模块提供给外部使用，是否是最小化的服务，是否可以独立完成工作。
很显然，想象一下，假设我有一个shop服务，想要注册用户、登录用户、获取用户信息等操作，只需要提供user rpc服务，就可以实现这个功能。而article rpc相关内容对于shop服务是多余的。
当然，你页可以把user rpc再拆分成 user rpc、login rpc、register rpc，如果你不嫌rpc服务太多，代码太分散的话。
