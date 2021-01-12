package serivces

//构造初始化交换器、队列



//服务接口
type Service interface {
	Pub()
	Sub()
}

type postAnswer struct {
	url string
}

func (s *postAnswer) Pub()  {

}