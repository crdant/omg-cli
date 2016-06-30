package file_server 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Diego struct {

	/*FileServer - Descr: Address of interface on which to serve files Default: 0.0.0.0:8080
*/
	FileServer *FileServer `yaml:"file_server,omitempty"`

	/*Ssl - Descr: when connecting over https, ignore bad ssl certificates Default: false
*/
	Ssl *Ssl `yaml:"ssl,omitempty"`

}