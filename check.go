package echo_health

type Check interface {
	CheckLiveness() error
	CheckReadiness() error
}

type CheckFunc func() error

func (cf CheckFunc) CheckReadiness() error {
	return cf()
}

func (cf CheckFunc) CheckLiveness() error {
	return cf()
}

type LivenessCheckFunc func() error

func (cf LivenessCheckFunc) CheckReadiness() error {
	return nil
}

func (cf LivenessCheckFunc) CheckLiveness() error {
	return cf()
}

type ReadinessCheckFunc func() error

func (cf ReadinessCheckFunc) CheckReadiness() error {
	return cf()
}

func (cf ReadinessCheckFunc) CheckLiveness() error {
	return nil
}
