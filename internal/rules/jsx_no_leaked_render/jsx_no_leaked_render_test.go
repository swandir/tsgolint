package jsx_no_leaked_render

import (
	"testing"

	"github.com/typescript-eslint/tsgolint/internal/rule_tester"
	"github.com/typescript-eslint/tsgolint/internal/rules/fixtures"
)

func TestJsxNoLeakedRenderRule(t *testing.T) {
	t.Parallel()

	rule_tester.RunRuleTester(fixtures.GetRootDir(), "tsconfig.minimal.json", t, &JsxNoLeakedRenderRule, []rule_tester.ValidTestCase{
		{
			Code: `
const a = <>{!(0) && <Foo />}</>;
const b = <>{!(NaN) && <Foo />}</>;
			`,
			Tsx: true,
		},
		{
			Code: `
const a = <>{!!(0) && <Foo />}</>;
const b = <>{!!(NaN) && <Foo />}</>;
			`,
			Tsx: true,
		},
		{
			Code: `
const a = <>{!!!(0) && <Foo />}</>;
const b = <>{!!!(NaN) && <Foo />}</>;
			`,
			Tsx: true,
		},
		{
			Code: `
let x: number | undefined;
const a = <>{!x && <Foo />}</>;
			`,
			Tsx: true,
		},
		{
			Code: `
let x: number | undefined;
const y = 2;
const a = <>{!x ? !x && <Foo /> : y && <Bar />}</>;
			`,
			Tsx: true,
		},
		{
			Code: `
const x = -1
const y = -0
const z = 0
const w = 1
const a = <>{x && <Foo />}</>;
const b = <>{!y && <Foo />}</>;
const c = <>{!z && <Foo />}</>;
const d = <>{w && <Foo />}</>;
			`,
			Tsx: true,
		},
		{
			Code: `
const x = -1n
const y = -0n
const z = 0n
const w = 1n
const a = <>{x && <Foo />}</>;
const b = <>{!y && <Foo />}</>;
const c = <>{!z && <Foo />}</>;
const d = <>{w && <Foo />}</>;
			`,
			Tsx: true,
		},
		{
			Code: `
const foo = Math.random() > 0.5;
const bar = "bar";
const a = <div>{0 || bar}</div>
			`,
			Tsx: true,
		},
		{
			Code: `
const foo = Math.random() > 0.5;
const bar = "bar";
const a = <div>{foo || bar}</div>
			`,
			Tsx: true,
		},
		{
			Code: `
type AppProps = { foo: string; }
const App = ({ foo }: AppProps) => <div>{foo}</div>
			`,
			Tsx: true,
		},
		{
			Code: `
type AppProps = {
  items: string[];
}
const App = ({ items }: AppProps) => {
  return <div>There are {items.length} elements</div>
}
			`,
			Tsx: true,
		},
		{
			Code: `
type AppProps = {
  items: string[];
  count: number;
}
const App = ({ items, count }: AppProps) => {
  return <div>{!count && 'No results found'}</div>
}
			`,
			Tsx: true,
		},
		{
			Code: `
type ListProps = { items: string[]; }
const List = ({ items }: ListProps) => {
  return <div>{items.map(item => <div key={item}>{item}</div>)}</div>
}
type AppProps = { items: string[]; }
const App = ({ items }: AppProps) => {
  return <div>{!!items.length && <List items={items}/>}</div>
}
			`,
			Tsx: true,
		},
		{
			Code: `
type ListProps = { items: string[]; }
const List = ({ items }: ListProps) => {
  return <div>{items.map(item => <div key={item}>{item}</div>)}</div>
}
type AppProps = { items: string[]; }
const App = ({ items }: AppProps) => {
  return <div>{Boolean(items.length) && <List items={items}/>}</div>
}
			`,
			Tsx: true,
		},
		{
			Code: `
type ListProps = { items: string[]; }
const List = ({ items }: ListProps) => {
  return <div>{items.map(item => <div key={item}>{item}</div>)}</div>
}
type AppProps = { items: string[]; }
const App = ({ items }: AppProps) => {
  return <div>{items.length > 0 && <List items={items}/>}</div>
}
			`,
			Tsx: true,
		},
		{
			Code: `
type ListProps = { items: string[]; }
const List = ({ items }: ListProps) => {
  return <div>{items.map(item => <div key={item}>{item}</div>)}</div>
}
type AppProps = { items: string[]; }
const App = ({ items }: AppProps) => {
  return <div>{items.length ? <List items={items}/> : null}</div>
}
			`,
			Tsx: true,
		},
		{
			Code: `
type ListProps = { items: string[]; }
const List = ({ items }: ListProps) => {
  return <div>{items.map(item => <div key={item}>{item}</div>)}</div>
}
type AppProps = { items: string[]; count: number; }
const App = ({ items, count }: AppProps) => {
  return <div>{count ? <List items={items}/> : null}</div>
}
			`,
			Tsx: true,
		},
		{
			Code: `
type ListProps = { items: string[]; }
const List = ({ items }: ListProps) => {
  return <div>{items.map(item => <div key={item}>{item}</div>)}</div>
}
type AppProps = { items: string[]; count: number; }
const App = ({ items, count }: AppProps) => {
  return <div>{!!count && <List items={items}/>}</div>
}
			`,
			Tsx: true,
		},
		{
			Code: `
const alwaysTruthy = true;
const alwaysFalsy = false;
const App = () => {
  return (
    <div>
      {alwaysTruthy && <div />}
      {alwaysFalsy && <div />}
    </div>
  )
}
			`,
			Tsx: true,
		},
		{
			Code: `
declare const direction: "down" | "up" | "";
const App = () => {
  return <div>{direction ? (direction === "down" ? "v" : "^") : ""}</div>
}
			`,
			Tsx: true,
		},
		{
			Code: `
const a = <>
  {0 ? <Foo /> : null}
  {'0' && <Foo />}
  {NaN ? <Foo /> : null}
</>
			`,
			Tsx: true,
		},
		{
			Code: `
const foo = Math.random() > 0.5;
const bar = 0;
function App() {
  return (
    <button
      type="button"
      disabled={foo && bar === 0}
      onClick={() => {}}
    />
  );
}
			`,
			Tsx: true,
		},
		{
			Code: `
const someCondition = JSON.parse("true") as boolean;
declare const val1: any;
declare const val2: any;
const SomeComponent = (_p: any) => <div />;
const a = <>{!!someCondition && (<SomeComponent prop1={val1} prop2={val2} />)}</>
			`,
			Tsx: true,
		},
		{
			Code: `
const someCondition = JSON.parse("") as any;
declare const val1: any;
declare const val2: any;
const SomeComponent = (_p: any) => <div />;
const a = <>{!!someCondition && (<SomeComponent prop1={val1} prop2={val2} />)}</>
			`,
			Tsx: true,
		},
		{
			Code: `
const someCondition = JSON.parse("") as unknown;
declare const val1: any;
declare const val2: any;
const SomeComponent = (_p: any) => <div />;
const a = <>{!!someCondition && (<SomeComponent prop1={val1} prop2={val2} />)}</>
			`,
			Tsx: true,
		},
		{
			Code: `
const someCondition = 0
declare const val1: any;
declare const val2: any;
const SomeComponent = (_p: any) => <div />;
const a = <>{!!someCondition && (<SomeComponent prop1={val1} prop2={val2} />)}</>
			`,
			Tsx: true,
		},
		{
			Code: `
const someCondition = 1
declare const val1: any;
declare const val2: any;
const SomeComponent = (_p: any) => <div />;
const a = <>{!!someCondition && (<SomeComponent prop1={val1} prop2={val2} />)}</>
			`,
			Tsx: true,
		},
		{
			Code: `
const someCondition = 0;
declare const val1: any;
declare const val2: any;
const SomeComponent = (_p: any) => <div />;
const a = <>{!!someCondition ? (<SomeComponent prop1={val1} prop2={val2} />) : someCondition ? null : <div />}</>
			`,
			Tsx: true,
		},
		{
			Code: `
const SomeComponent = () => <div />;
const App = ({
  someCondition,
}: {
  someCondition?: number | undefined;
}) => {
  return (
    <>
      {someCondition
        ? someCondition
        : <SomeComponent />
      }
    </>
  )
}
			`,
			Tsx: true,
		},
		{
			Code: `
const someCondition = 0;
declare const val1: any;
declare const val2: any;
const SomeComponent = (_p: any) => <div />;
const App = () => {
  return (
    <>
      {!!someCondition
        ? (
        <SomeComponent
          prop1={val1}
          prop2={val2}
        />)
        : someCondition ? "aaa"
        : someCondition && someCondition
        ? <div />
        : null
      }
    </>
  )
}
			`,
			Tsx: true,
		},
		{
			Code: `
const someCondition = true
declare const val1: any;
declare const val2: any;
const SomeComponent = (_p: any) => <div />;
const App = () => {
  return (
    <>
      {someCondition
        ? (
        <SomeComponent
          prop1={val1}
          prop2={val2}
        />)
        : "else"
      }
    </>
  )
}
			`,
			Tsx: true,
		},
		{
			Code: `
const someCondition = 0;
declare const val1: any;
declare const val2: any;
const SomeComponent = (_p: any) => <div />;
const App = () => {
  return (
    <>
      {!!someCondition
        ? (
        <SomeComponent
          prop1={val1}
          prop2={val2}
        />)
        : someCondition
        ? someCondition
        : <div />
      }
    </>
  )
}
			`,
			Tsx: true,
		},
		{
			Code: `
const SomeComponent = () => <div />;
const someFunction = (input: unknown): 10 => 10
const App = ({ someCondition }: { someCondition?: number | undefined }) => {
  return <>{someCondition ? someFunction(someCondition) : <SomeComponent />}</>;
};
			`,
			Tsx: true,
		},
		{
			Code: `
const SomeComponent = () => <div />;
const App = ({
  someCondition,
}:{
  someCondition?: boolean | undefined;
}) => {
  return <>{someCondition && <SomeComponent />}</>;
}
			`,
			Tsx: true,
		},
		{
			Code: `
type AppProps<T> = {
  someFunction: (data: T) => unknown;
};
function App<T>({ someFunction }: AppProps<T>) {
  return <>{someFunction && someFunction(1 as T)}</>;
}
			`,
			Tsx: true,
		},
		{
			Code: `
function App<T extends string>({ value }: { value: T }) {
  return <>{value && <Foo />}</>;
}
			`,
			Tsx: true,
		},
		{
			Code: `
function App<T extends object>({ value }: { value: T }) {
  return <>{value && <Foo />}</>;
}
			`,
			Tsx: true,
		},
		{
			Code: `
type Name = string & { readonly __brand: unique symbol };
function App({ name }: { name: Name }) {
  return <>{name && <Foo />}</>;
}
			`,
			Tsx: true,
		},
		{
			Code: `
const someCondition = JSON.parse("") as any;
declare const val1: any;
declare const val2: any;
const SomeComponent = (_p: any) => <div />;
const App = () => {
  return (
    <>
      {someCondition && (
        <SomeComponent
          prop1={val1}
          prop2={val2}
        />
      )}
    </>
  )
}
			`,
			Tsx: true,
		},
		{
			Code: `
function App() {
  const a = {} as {};
  const b = {} as {} | null;
  const c = {} as {} | undefined;
  return (
    <>
      <>{a && <div />}</>
      <>{b && <div />}</>
      <>{c && <div />}</>
    </>
  );
}
			`,
			Tsx: true,
		},
		{
			Code: `
declare function getData(): { id: number; name: string }[] | undefined;
function List({ items }: { items: { id: number; name: string }[] }) {
  return <div>{items.map(item => <div key={item.id}>{item.name}</div>)}</div>;
}
type Item = { id: number; name: string };
function App() {
  let data: Item[] | undefined = getData();
  return (
    <div>
      {data && <List items={data} />}
    </div>
  );
}
			`,
			Tsx: true,
		},
		{
			Code: `
export const MyComponent = ({ isVisible1 }: { isVisible1: boolean }) => {
  const isVisible2 = true;
  const isVisible3 = 1 > 2;
  return (
    <>
      {isVisible1 && <div />}
      {isVisible2 && <div />}
      {isVisible3 && <div />}
    </>
  );
};
			`,
			Tsx: true,
		},
		{
			Code: `
enum A {
  Foo = 'foo',
  Bar = 1,
}
const App = ({ value }: { value: A }) => {
  return <>{value && <div>{value}</div>}</>;
};
			`,
			Tsx: true,
		},
		{
			Code: `
const someString = "";
const a = <>{someString && <Something />}</>;
			`,
			Tsx: true,
		},
		{
			Code: `
const anyString = Math.random() > 0.5 ? "" : "foo";
const a = <>{anyString && <Something />}</>;
			`,
			Tsx: true,
		},
		{
			Code: `
const a = <>{"" && <Something />}</>;
			`,
			Tsx: true,
		},
	}, []rule_tester.InvalidTestCase{
		{
			Code: `
const a = <>{0 && <Foo />}</>;
const b = <>{NaN && <Foo />}</>;
			`,
			Tsx: true,
			Errors: []rule_tester.InvalidTestCaseError{
				{MessageId: "default"},
				{MessageId: "default"},
			},
		},
		{
			Code: `
const a = <>{(0) && <Foo />}</>;
const b = <>{(NaN) && <Foo />}</>;
			`,
			Tsx: true,
			Errors: []rule_tester.InvalidTestCaseError{
				{MessageId: "default"},
				{MessageId: "default"},
			},
		},
		{
			Code: `
const x = -1
const y = -0
const z = 0
const w = 1
const a = <>{x && <Foo />}</>;
const b = <>{y && <Foo />}</>;
const c = <>{z && <Foo />}</>;
const d = <>{w && <Foo />}</>;
			`,
			Tsx: true,
			Errors: []rule_tester.InvalidTestCaseError{
				{MessageId: "default"},
				{MessageId: "default"},
			},
		},
		{
			Code: `
const x = -1n
const y = -0n
const z = 0n
const w = 1n
const a = <>{x && <Foo />}</>;
const b = <>{y && <Foo />}</>;
const c = <>{z && <Foo />}</>;
const d = <>{w && <Foo />}</>;
			`,
			Tsx: true,
			Errors: []rule_tester.InvalidTestCaseError{
				{MessageId: "default"},
				{MessageId: "default"},
			},
		},
		{
			Code: `
const someCondition = JSON.parse("") as unknown;
declare const val1: any;
declare const val2: any;
const SomeComponent = (_p: any) => <div />;
const a = <>{someCondition && <SomeComponent prop1={val1} prop2={val2} />}</>;
			`,
			Tsx: true,
			Errors: []rule_tester.InvalidTestCaseError{
				{MessageId: "default"},
			},
		},
		{
			Code: `
const someCondition = 0;
declare const val1: any;
declare const val2: any;
const SomeComponent = (_p: any) => <div />;
const a = <>{someCondition && <SomeComponent prop1={val1} prop2={val2} />}</>;
			`,
			Tsx: true,
			Errors: []rule_tester.InvalidTestCaseError{
				{MessageId: "default"},
			},
		},
		{
			Code: `
const someCondition = -0;
declare const val1: any;
declare const val2: any;
const SomeComponent = (_p: any) => <div />;
const a = <>{someCondition && <SomeComponent prop1={val1} prop2={val2} />}</>;
			`,
			Tsx: true,
			Errors: []rule_tester.InvalidTestCaseError{
				{MessageId: "default"},
			},
		},
		{
			Code: `
const someCondition = 0;
declare const val1: any;
declare const val2: any;
const SomeComponent = (_p: any) => <div />;
const a = <>{!!someCondition ? <SomeComponent prop1={val1} prop2={val2} /> : someCondition && <div />}</>;
			`,
			Tsx: true,
			Errors: []rule_tester.InvalidTestCaseError{
				{MessageId: "default"},
			},
		},
		{
			Code: `
const App = ({
  someCondition,
}: {
  someCondition: number | undefined;
}) => {
  return (
    <>
      {someCondition && <Foo />}
    </>
  )
}
			`,
			Tsx: true,
			Errors: []rule_tester.InvalidTestCaseError{
				{MessageId: "default"},
			},
		},
		{
			Code: `
const SomeComponent = () => <div />;
const App = ({
  someCondition,
}:{
  someCondition?: number | undefined;
}) => {
  return (
    <>
      {someCondition && <SomeComponent />}
    </>
  )
}
			`,
			Tsx: true,
			Errors: []rule_tester.InvalidTestCaseError{
				{MessageId: "default"},
			},
		},
		{
			Code: `
enum A {
  Foo = 'foo',
  Bar = 0,
}
const App = ({ value }: { value: A }) => {
  return <>{value && <div>{value}</div>}</>;
};
			`,
			Tsx: true,
			Errors: []rule_tester.InvalidTestCaseError{
				{MessageId: "default"},
			},
		},
		{
			Code: `
const App = ({ value }: { value: number }) => {
  const foo = value && <Foo />;
  return <>{foo}</>;
};
			`,
			Tsx: true,
			Errors: []rule_tester.InvalidTestCaseError{
				{MessageId: "default"},
			},
		},
		{
			Code: `
function App<T>({ value }: { value: T }) {
  return <>{value && <Foo />}</>;
}
			`,
			Tsx: true,
			Errors: []rule_tester.InvalidTestCaseError{
				{MessageId: "default"},
			},
		},
		{
			Code: `
function App<T extends number>({ value }: { value: T }) {
  return <>{value && <Foo />}</>;
}
			`,
			Tsx: true,
			Errors: []rule_tester.InvalidTestCaseError{
				{MessageId: "default"},
			},
		},
		{
			Code: `
type Count = number & { readonly __brand: unique symbol };
function App({ count }: { count: Count }) {
  return <>{count && <Foo />}</>;
}
			`,
			Tsx: true,
			Errors: []rule_tester.InvalidTestCaseError{
				{MessageId: "default"},
			},
		},
	})
}
