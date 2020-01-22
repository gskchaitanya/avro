package avroregistry_test

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"gopkg.in/retry.v1"

	"github.com/heetch/avro"
	"github.com/heetch/avro/avroregistry"
)

func TestRegister(t *testing.T) {
	c := qt.New(t)
	defer c.Done()
	r := newTestRegistry(c)

	type R struct {
		X int
	}
	ctx := context.Background()
	id, err := r.Register(ctx, schemaOf(nil, R{}))
	c.Assert(err, qt.Equals, nil)

	id1, err := r.IDForSchema(ctx, schemaOf(nil, R{}))
	c.Assert(err, qt.Equals, nil)
	c.Assert(id1, qt.Equals, id)
}

func TestSchemaCompatibility(t *testing.T) {
	c := qt.New(t)
	defer c.Done()
	r := newTestRegistry(c)
	ctx := context.Background()
	err := r.SetCompatibility(ctx, avro.BackwardTransitive)
	c.Assert(err, qt.Equals, nil)

	type R struct {
		X int
	}
	_, err = r.Register(ctx, schemaOf(nil, R{}))
	c.Assert(err, qt.Equals, nil)

	// Try to register an incompatible type.
	type R1 struct {
		X string
	}
	names := new(avro.Names).RenameType(R1{}, "R")
	_, err = r.Register(ctx, schemaOf(names, R1{}))
	c.Assert(err, qt.ErrorMatches, `Avro registry error \(code 409\): Schema being registered is incompatible with an earlier schema`)

	// Check that we can't rename the schema.
	_, err = r.Register(ctx, schemaOf(nil, R1{}))
	c.Assert(err, qt.ErrorMatches, `Avro registry error \(code 409\): Schema being registered is incompatible with an earlier schema`)

	// Check that we can change the field to a compatible union.
	type R2 struct {
		X *int
	}
	names = new(avro.Names).RenameType(R2{}, "R")
	_, err = r.Register(ctx, schemaOf(names, R2{}))
	c.Assert(err, qt.Equals, nil)

	// Check that we can't change it back again.
	type R3 struct {
		X int
		Y string
	}
	names = new(avro.Names).RenameType(R3{}, "R")
	_, err = r.Register(ctx, schemaOf(names, R3{}))
	c.Assert(err, qt.ErrorMatches, `Avro registry error \(code 409\): Schema being registered is incompatible with an earlier schema`)
}

func TestSingleCodec(t *testing.T) {
	c := qt.New(t)
	defer c.Done()
	r := newTestRegistry(c)
	ctx := context.Background()
	err := r.SetCompatibility(ctx, avro.BackwardTransitive)
	c.Assert(err, qt.Equals, nil)
	type R struct {
		X int
	}
	id1, err := r.Register(ctx, schemaOf(nil, R{}))
	c.Assert(err, qt.Equals, nil)

	type R1 struct {
		X int
		Y int
	}
	names := new(avro.Names).RenameType(R1{}, "R")
	id2, err := r.Register(ctx, schemaOf(names, R1{}))
	c.Assert(err, qt.Equals, nil)
	c.Assert(id2, qt.Not(qt.Equals), id1)

	enc := avro.NewSingleEncoder(r, names)
	data1, err := enc.Marshal(ctx, R{10})
	c.Assert(err, qt.Equals, nil)
	c.Assert(data1[0], qt.Equals, byte(0))
	c.Assert(int64(binary.BigEndian.Uint32(data1[1:5])), qt.Equals, id1)
	c.Assert(data1[5:], qt.DeepEquals, []byte{20})

	data2, err := enc.Marshal(ctx, R1{11, 30})
	c.Assert(err, qt.Equals, nil)
	c.Assert(data2[0], qt.Equals, byte(0))
	c.Assert(int64(binary.BigEndian.Uint32(data2[1:5])), qt.Equals, id2)
	c.Assert(data2[5:], qt.DeepEquals, []byte{22, 60})

	// Check that it round-trips back through the decoder OK.
	dec := avro.NewSingleDecoder(r, names)
	var x1 R
	_, err = dec.Unmarshal(ctx, data1, &x1)
	c.Assert(err, qt.Equals, nil)
	c.Assert(x1, qt.Equals, R{10})

	var x2 R1
	_, err = dec.Unmarshal(ctx, data2, &x2)
	c.Assert(err, qt.Equals, nil)
	c.Assert(x2, qt.Equals, R1{11, 30})
}

func TestRetryOnError(t *testing.T) {
	c := qt.New(t)
	defer c.Done()
	registry, err := avroregistry.New(avroregistry.Params{
		ServerURL: "http://0.1.2.3",
		Subject:   "subject",
		RetryStrategy: retry.LimitCount(4, retry.Regular{
			Total: time.Second,
			Delay: 10 * time.Millisecond,
		}),
	})
	c.Assert(err, qt.Equals, nil)
	t0 := time.Now()
	err = registry.SetCompatibility(context.Background(), avro.BackwardTransitive)
	c.Assert(err, qt.ErrorMatches, `Put "?http://0.1.2.3/config/subject"?: dial tcp 0.1.2.3:80: connect: .*`)
	if d := time.Since(t0); d < 30*time.Millisecond {
		c.Errorf("retry duration too small, want >=30ms got %v", d)
	}
}

func TestCanceledRetry(t *testing.T) {
	c := qt.New(t)
	defer c.Done()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()
	registry, err := avroregistry.New(avroregistry.Params{
		ServerURL: "http://0.1.2.3",
		Subject:   "subject",
	})
	c.Assert(err, qt.Equals, nil)
	t0 := time.Now()
	err = registry.SetCompatibility(ctx, avro.BackwardTransitive)
	c.Assert(err, qt.ErrorMatches, `context canceled`)
	if d := time.Since(t0); d > 500*time.Millisecond {
		c.Errorf("retry duration too large, want ~30ms got %v", d)
	}
}

func schemaOf(names *avro.Names, x interface{}) string {
	if names == nil {
		names = new(avro.Names)
	}
	t, err := names.TypeOf(x)
	if err != nil {
		panic(err)
	}
	return t.String()
}

func newTestRegistry(c *qt.C) *avroregistry.Registry {
	ctx := context.Background()
	serverURL := os.Getenv("AVRO_REGISTRY_URL")
	if serverURL == "" {
		c.Skip("no AVRO_REGISTRY_URL configured")
	}
	subject := randomString()
	registry, err := avroregistry.New(avroregistry.Params{
		ServerURL:     serverURL,
		Subject:       subject,
		RetryStrategy: noRetry,
	})
	c.Assert(err, qt.Equals, nil)
	c.Defer(func() {
		err := registry.DeleteSubject(ctx)
		c.Check(err, qt.Equals, nil)
	})
	return registry
}

var noRetry = retry.Regular{}

func randomString() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("test-%x", buf)
}
