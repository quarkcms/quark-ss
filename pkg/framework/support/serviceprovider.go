package serviceprovider

type Provider interface {
	Register()
	Boot()
}
