package jq

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)

type ValidationError struct {
	s string
}

func (e *ValidationError) Error() string {
	return e.s
}

var (
	ExecTimeoutError   = errors.New("jq execution was timeout")
	ExecCancelledError = errors.New("jq execution was cancelled")
)

type JQ struct {
	J string  `json:"j"`
	Q string  `json:"q"`
	O []JQOpt `json:"o"`
}

type JQOpt struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func (j *JQ) Opts() []string {
	opts := []string{}
	for _, opt := range j.O {
		if opt.Enabled {
			opts = append(opts, fmt.Sprintf("--%s", opt.Name))
		}
	}

	return opts
}

func (j *JQ) EvalStream(
	ctx context.Context,
	r io.Reader,
	w io.Writer,
	e io.Writer,
) error {
	// if err := j.Validate(); err != nil {
	// 	return err
	// }
	//
	// var outbuf, errbuf bytes.Buffer
	// buf := new(bytes.Buffer)
	// buf.ReadFrom(r)
	// s := buf.String()
	// log.Println("json: ", len(s))

	opts := j.Opts()
	opts = append(opts, strings.TrimSpace(j.Q))
	log.Println(opts)
	cmd := exec.CommandContext(ctx, Path, opts...)
	cmd.Stdin = r
	cmd.Env = make([]string, 0)
	cmd.Stdout = w
	cmd.Stderr = e
	// cmd.Stdout = &outbuf
	// cmd.Stderr = &errbuf
	// fmt.Fprint(w, "hello")

	err := cmd.Run()
	if err != nil {
		ctxErr := ctx.Err()

		if ctxErr == context.DeadlineExceeded {
			return ExecTimeoutError
		}
		if ctxErr == context.Canceled {
			return ExecCancelledError
		}
	}

	// stdout := outbuf.String()
	// stderr := errbuf.String()
    //
	// log.Println(stdout)
	// log.Println(stderr)
	return err
}

type logWriter struct{ *log.Logger }

func (w logWriter) Write(b []byte) (int, error) {
	w.Printf("writing: \n")
	w.Printf("%s", b)
	return len(b), nil
}

func (j *JQ) Eval(
	ctx context.Context,
	w io.Writer,
	e io.Writer,
) error {
	if err := j.Validate(); err != nil {
		return err
	}

	opts := j.Opts()
	opts = append(opts, j.Q)
	cmd := exec.CommandContext(ctx, Path, opts...)
	cmd.Stdin = bytes.NewBufferString(j.J)
	cmd.Env = make([]string, 0)
	cmd.Stdout = w
	cmd.Stderr = e

	err := cmd.Run()
	if err != nil {
		ctxErr := ctx.Err()

		if ctxErr == context.DeadlineExceeded {
			return ExecTimeoutError
		}
		if ctxErr == context.Canceled {
			return ExecCancelledError
		}
	}

	return err
}

func (j *JQ) ValidateFilter() error {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var b bytes.Buffer
	e := bufio.NewWriter(&b)
	cmd := exec.CommandContext(ctx, Path, j.Q)
	cmd.Stdin = bytes.NewBufferString("")
	cmd.Env = make([]string, 0)

	cmd.Stderr = e

	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
		e.Flush()
		log.Println(b.String())
	}
	return err
}

func (j *JQ) Validate() error {
	errMsgs := []string{}

	if j.Q == "" {
		errMsgs = append(errMsgs, "missing filter")
	}

	if j.J == "" {
		errMsgs = append(errMsgs, "missing JSON")
	}

	if len(errMsgs) > 0 {
		return &ValidationError{fmt.Sprintf("invalid input: %s", strings.Join(errMsgs, " and "))}
	}

	return nil
}

func (j JQ) String() string {
	return fmt.Sprintf("j=%s, q=%s, o=%v", j.J, j.Q, j.Opts())
}
