package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// runBindDestroyer create intermediate communication link which breaks on chDone channel close.
func runBindDestroyer(chIn In, chDone In) Out {
	chOut := make(Bi)

	go func() {
		defer close(chOut)

		for {
			// when we handle several channels and all of them are ready, runtime choose random
			// increase priority of done channel
			select {
			case <-chDone:
				return
			default:
			}

			select {
			case val, ok := <-chIn:
				if !ok {
					return
				}

				select {
				case <-chDone:
					return
				default:
				}

				select {
				case <-chDone:
					return
				case chOut <- val:
				}
			case <-chDone:
				return
			}
		}
	}()

	return chOut
}

// ExecutePipeline combine execution pipe, binds prev chOut to next chIn,
// place between them special link, which breaks on chDone close.
func ExecutePipeline(chIn In, chDone In, stages ...Stage) Out {
	chOut := chIn

	for _, s := range stages {
		chOut = runBindDestroyer(chOut, chDone)
		chOut = s(chOut)
	}

	return chOut
}
