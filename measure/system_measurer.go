// Create package in v.1.0.0
// measure package define measurer struct about system, msg(문자 서비스) usage, etc ...
// system_measurer.go define system measurer about cpu or memory usage, disk remain capacity, etc ...

package measure

// systemMeasurer is struct that measure value about system
type systemMeasurer struct {}

// SystemMeasurer function return systemMeasurer ptr instance with initializing
func SystemMeasurer() *systemMeasurer {
	return &systemMeasurer{}
}
