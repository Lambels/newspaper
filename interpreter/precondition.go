package interpreter

type Precondition interface {
    EvalElement

    //TODO: add some Context
    Advance(Context) error
}
