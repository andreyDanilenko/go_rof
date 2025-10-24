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
	stageOut := stage(current)

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case value, ok := <-stageOut:
				if !ok {
					return
				}
				out <- value
			}
		}
	}()

	return out
}
