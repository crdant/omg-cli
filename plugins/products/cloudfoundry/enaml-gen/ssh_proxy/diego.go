package ssh_proxy 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Diego struct {

	/*SshProxy - Descr: address at which to serve debug info Default: 0.0.0.0:17016
*/
	SshProxy *SshProxy `yaml:"ssh_proxy,omitempty"`

	/*Ssl - Descr: when connecting over https, ignore bad ssl certificates Default: false
*/
	Ssl *Ssl `yaml:"ssl,omitempty"`

}