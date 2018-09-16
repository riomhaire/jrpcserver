package none

type DefaultServiceRegistry struct {
}

func NewDefaultServiceRegistry() (r *DefaultServiceRegistry) {
	r = &DefaultServiceRegistry{}
	return
}

func (r *DefaultServiceRegistry) Register() error {
	return nil

}

func (r *DefaultServiceRegistry) Deregister() error {
	return nil

}
