package pipies

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koss-null/funcfrog/internal/primitive/pointer"
)

type testT string

func Test_Predicates(t *testing.T) {
	t.Parallel()

	t.Run("NotNil", func(t *testing.T) {
		require.True(t, NotNil(pointer.To(1)))
		require.False(t, NotNil[int](nil))
		require.False(t, NotNil[int](nil))
		require.False(t, NotNil[testT](nil))
		var empty *testT
		require.False(t, NotNil(empty))
		empty = nil
		require.False(t, NotNil(empty))
		var e any = empty
		require.True(t, NotNil[any](&e))
	})
	t.Run("IsNil", func(t *testing.T) {
		require.False(t, IsNil(pointer.To(1)))
		require.True(t, IsNil[int](nil))
		require.True(t, IsNil[int](nil))
		require.True(t, IsNil[testT](nil))
		var empty *testT
		require.True(t, IsNil(empty))
		empty = nil
		require.True(t, IsNil(empty))
		var e any = empty
		require.False(t, IsNil[any](&e))
	})
	t.Run("NotZero", func(t *testing.T) {
		require.True(t, NotZero(pointer.To(1)))
		require.False(t, NotZero(pointer.To(0)))
		require.False(t, NotZero(pointer.To[float32](0.0)))
		require.False(t, NotZero(&struct{ a, b, c int }{}))
		require.True(t, NotZero(&struct{ a, b, c int }{1, 0, 0}))
	})
}

func Test_PredicateBuilders(t *testing.T) {
	t.Parallel()

	t.Run("Eq", func(t *testing.T) {
		eq5 := Eq(5)
		require.True(t, eq5(pointer.To(5)))
		require.False(t, eq5(pointer.To(4)))
		eqS := Eq("test")
		require.True(t, eqS(pointer.To("test")))
		require.False(t, eqS(pointer.To("sett")))
	})

	t.Run("NotEq", func(t *testing.T) {
		neq5 := NotEq(5)
		require.False(t, neq5(pointer.To(5)))
		require.True(t, neq5(pointer.To(4)))
		neqS := NotEq("test")
		require.False(t, neqS(pointer.To("test")))
		require.True(t, neqS(pointer.To("sett")))
	})

	t.Run("LessThan", func(t *testing.T) {
		lt5 := LessThan(5)
		require.True(t, lt5(pointer.To(4)))
		require.False(t, lt5(pointer.To(5)))
		require.False(t, lt5(pointer.To(6)))
		ltf5 := LessThan(5.0)
		require.True(t, ltf5(pointer.To(4.999)))
		require.False(t, ltf5(pointer.To(float64(int(5)))))
		require.False(t, ltf5(pointer.To(5.01)))
	})
}

func Test_Comparator(t *testing.T) {
	require.True(t, Less(pointer.To(4), pointer.To(5)))
	require.False(t, Less(pointer.To(5), pointer.To(5)))
	require.False(t, Less(pointer.To(6), pointer.To(5)))
	require.True(t, Less(pointer.To(4.999), pointer.To(5.0)))
	require.False(t, Less(pointer.To(float64(int(5))), pointer.To(5.0)))
	require.False(t, Less(pointer.To(5.01), pointer.To(5.0)))
}

func Test_Accum(t *testing.T) {
	require.Equal(t, Sum(pointer.To(10), pointer.To(20)), 30)
	require.Equal(t, Sum(pointer.To(10.0), pointer.To(20.0)), 30.0)
}
