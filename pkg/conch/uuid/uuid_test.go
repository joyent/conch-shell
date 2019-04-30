package uuid_test

import (
	"encoding/json"
	"testing"

	"github.com/joyent/conch-shell/pkg/conch/uuid"
	"github.com/nbio/st"
)

func TestUUID(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		u := uuid.New()
		st.Expect(t, u, uuid.UUID{})
		st.Expect(t, u.IsZero(), true)
		st.Expect(t, u.Equal(u), true)
	})

	t.Run("V4", func(t *testing.T) {
		u2 := uuid.NewV4()
		st.Reject(t, u2, uuid.UUID{})
		st.Expect(t, u2.IsZero(), false)

		u := uuid.New()
		st.Expect(t, u.Equal(u2), false)
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		u := uuid.NewV4()
		j, err := json.Marshal(u)

		st.Expect(t, err, nil)
		st.Expect(t, string(j), "\""+u.String()+"\"")
	})

	t.Run("UnmarshalJSON Round trip", func(t *testing.T) {
		var u2 uuid.UUID
		u := uuid.NewV4()

		j, err := json.Marshal(u)
		st.Expect(t, err, nil)
		st.Expect(t, string(j), "\""+u.String()+"\"")

		err = json.Unmarshal(j, &u2)
		st.Expect(t, err, nil)

		st.Expect(t, u, u2)
		st.Expect(t, u.String(), u2.String())
		st.Expect(t, u2.IsZero(), false)
		st.Expect(t, u.Equal(u2), true)
	})

}
