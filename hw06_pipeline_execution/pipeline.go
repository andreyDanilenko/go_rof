package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	current := in

	for _, stage := range stages {
		current = LayerStage(stage, current, done)
	}

	return current
}

func LayerStage(stage Stage, current In, done In) Out {
	out := make(Bi)
	layerCurrent := make(Bi)

	go func() {
		defer close(layerCurrent)
		for v := range current {
			select {
			case <-done:
				return
			case layerCurrent <- v:
			}
		}
	}()

	stageOut := stage(layerCurrent)

	go func() {
		defer close(out)
		for v := range stageOut {
			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}()

	return out
}
