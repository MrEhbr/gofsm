package fsm

var fsmTemplate = `package {{.Package.Name}}
// DO NOT EDIT!
// This code is generated with http://github.com/MrEhbr/gofsm tool

{{ if not .Options.DisableGoGenerate}}
//{{"go:generate"}} gofsm gen -s {{.Options.Struct}} -f {{.Options.StateField}} -o {{ base .Options.OutputFile }}
{{- if .Options.TransitionsFile}} -t {{ path_join (rel (dir .Options.OutputFile) (dir .Options.TransitionsFile)) (base .Options.TransitionsFile) }} {{- end }}
{{- if .Options.ActionGraphOutputFile}} -a {{ path_join (rel (dir .Options.OutputFile) (dir .Options.ActionGraphOutputFile)) (base .Options.ActionGraphOutputFile) }} {{- end }}
{{- end }}

type (
	// {{.Struct.Name}}Transition is a state transition and all data are literal values that simplifies FSM usage and make it generic.
	{{.Struct.Name}}Transition struct {
		Event string
		From {{.Struct.StateType}}
		To {{.Struct.StateType}}
		BeforeActions []string
		Actions []string
	}
	// {{.Struct.Name}}Handle handles transitions action
	{{.Struct.Name}}HandleAction func(action string, fromState, toState {{.Struct.StateType}}, obj *{{.Struct.Name}}) error
	// Save state to external storage
	{{.Struct.Name}}PersistState func(obj *{{.Struct.Name}}, state {{.Struct.StateType}}) error
	// {{.Struct.Name}}StateMachine is a FSM that can handle transitions of a lot of objects. eventHandler and transitions are configured before use them.
	{{.Struct.Name}}StateMachine struct {
		transitions []{{.Struct.Name}}Transition
		actionHandler {{.Struct.Name}}HandleAction
		persister {{.Struct.Name}}PersistState
	}
)

var (
	Err{{.Struct.Name}}FsmAction = errors.New("{{.Struct.Name}}StateMachine action error")
	Err{{.Struct.Name}}FsmBeforeAction = errors.New("{{.Struct.Name}}StateMachine before action error")
	// Err{{.Struct.Name}}Skip indicates that further processing not need
	// used in before_actions
	Err{{.Struct.Name}}FsmSkip = errors.New("skip")
)

type Option func(*{{.Struct.Name}}StateMachine)

func WithActionHandler(h {{.Struct.Name}}HandleAction) Option {
	return func(fsm *{{.Struct.Name}}StateMachine) {
		fsm.actionHandler = h
	}
}

func WithPersiter(p {{.Struct.Name}}PersistState) Option {
	return func(fsm *{{.Struct.Name}}StateMachine) {
		fsm.persister = p
	}
}

func WithTransitions(tr []{{.Struct.Name}}Transition) Option {
	return func(fsm *{{.Struct.Name}}StateMachine) {
		fsm.transitions = tr
	}
}

// New{{.Struct.Name}}StateMachine creates a new state machine.
func New{{.Struct.Name}}StateMachine(opts ...Option) *{{.Struct.Name}}StateMachine {
	fsm := &{{.Struct.Name}}StateMachine{}
	for _, o := range opts {
		o(fsm)
	}

	return fsm
}

// ChangeState fires a event and if event succeeded then change state.
func (m *{{.Struct.Name}}StateMachine) ChangeState(event string, obj *{{.Struct.Name}}) error {
	trans, ok := m.findTransMatching(obj.{{.Struct.StateField}}, event)
	if !ok {
		return fmt.Errorf("cannot find transition for event [%s] when in state [%v]", event, obj.{{.Struct.StateField}})
	}

	if len(trans.BeforeActions) > 0 && m.actionHandler != nil {
		for _, action := range trans.BeforeActions {
			if err := m.actionHandler(action, trans.From, trans.To, obj); err != nil {
				if errors.Is(err, Err{{.Struct.Name}}Skip) {
					return nil
				}

				return fmt.Errorf("%w. action [%s] return error: %s", Err{{.Struct.Name}}FsmBeforeAction, action, err)
			}
		}
	}

	if m.persister != nil {
		if err := m.persister(obj, trans.To); err != nil {
			return err
		}
	}

	obj.{{.Struct.StateField}} = trans.To

	if len(trans.Actions) > 0 && m.actionHandler  != nil {
		var errs error
		for _, action := range trans.Actions {
			if err := m.actionHandler(action, trans.From, trans.To, obj); err != nil {
				errs = multierror.Append(errs, fmt.Errorf("%w. action [%s] return error: %s", Err{{.Struct.Name}}FsmAction, action, err))
			}
		}

		if errs != nil {
			return errs
		}
	}

	return nil
}

func (m *{{.Struct.Name}}StateMachine) Can(state {{.Struct.StateType}}, event string) bool {
	_, ok := m.findTransMatching(state, event)
	return ok
}

func (m *{{.Struct.Name}}StateMachine) FindTransitionForStates(from, to {{.Struct.StateType}}) ({{.Struct.Name}}Transition, bool) {
	for _, v := range m.transitions {
		if v.From == from && v.To == to {
			return v, true
		}
	}
	return {{.Struct.Name}}Transition{}, false
}

// findTransMatching gets corresponding transition according to current state and event.
func (m *{{.Struct.Name}}StateMachine) findTransMatching(fromState {{.Struct.StateType}}, event string) ({{.Struct.Name}}Transition, bool) {
	for _, v := range m.transitions {
		if v.From == fromState && v.Event == event {
			return v, true
		}
	}
	return {{.Struct.Name}}Transition{}, false
}

{{- if .Struct.Transitions }}
const (
	// {{ .Struct.Name }} state machine events
{{- range $val := .Transitions }}
{{- if .Event }}
	{{ $.Struct.Name }}Event{{to_camel .Event true}} = "{{.Event}}"
{{- end }}
{{- end }}

	// {{ $.Struct.Name }} state machine actions
{{- range $val := .Transitions.Actions }}
	{{ $.Struct.Name }}Action{{to_camel $val true}} = "{{$val}}"
{{- end }}
)

// {{ $.Struct.Name }}Transitions generated from {{ path_join (rel (dir .Options.OutputFile) (dir .Options.TransitionsFile)) (base .Options.TransitionsFile) }}
var {{ $.Struct.Name }}Transitions = []{{ $.Struct.Name }}Transition{
{{- range $val := .Transitions }}
	{{- range $from := .From }}
{
{{- if $val.Event }}
	Event: {{ $.Struct.Name }}Event{{to_camel $val.Event true}},
{{- end }}
	From: {{ $.Struct.FindValue $from }},
	To: {{ $.Struct.FindValue $val.To }},
	{{- if $val.BeforeActions}}
	BeforeActions: []string{
		{{- range $action := $val.BeforeActions }}
			{{ $.Struct.Name }}Action{{to_camel $action true}},
		{{- end }}
	},
	{{- end }}

	{{- if $val.Actions}}
	Actions: []string{
		{{- range $action := $val.Actions }}
			{{ $.Struct.Name }}Action{{to_camel $action true}},
		{{- end }}
	},
	{{- end }}
},
	{{- end }}
{{- end }}
}
{{- end }}
`
