# _hulk_client
Hulk Golang语言客户端

## How to this SDK?
> 当_hulk_client被引入时, 其会进行初始化操作。 在初始化时就会自动去获取配置数据. 因此需要将_hulk_client放在import靠前的位置(虽然不优雅，但目前尚未找到更好的解决方案)

以下是初始化时所必须的环境变量:

* HULK_ENDPOINT HulkAPI地址 例如: https://hulk.devexp.cn/api
* HULK_PROJECT_NAME 需要获取配置数据的项目名称
* HULK_PROJECT_VERSION 项目版本

### 兼容Dep

在使用dep管理vendor时，如果没有使用package的任何一个函数，dep不会自动将此package放入vendor之中。 因此_hulk_client提供了一个Run函数。 此函数没有任何'副作用'。就是为了满足dep而产生的
