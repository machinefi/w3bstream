package deploy

type Deployer interface {
	Name() string
	Write(cwd string) error
}
