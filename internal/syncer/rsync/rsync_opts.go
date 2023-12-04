package rsync

import (
	"errors"
	"fmt"
)

func BandwidthLimit(limit string) RsyncOpt {
	return func(r *Rsync) error {
		if len(limit) == 0 {
			return errors.New("empty limit supplied")
		}

		s := fmt.Sprintf(`--bwlimit=%s`, limit)
		r.args = append(r.args, s)
		return nil
	}
}

func RemoveSourceFiles() func(r *Rsync) error {
	return func(r *Rsync) error {
		r.args = append(r.args, "--remove-source-files")
		return nil
	}
}

func Exclude(exclude string) func(r *Rsync) error {
	return func(r *Rsync) error {
		if len(exclude) == 0 {
			return errors.New("empty exclude pattern supplied")
		}

		s := fmt.Sprintf(`--exclude="%s"`, exclude)
		r.args = append(r.args, s)
		return nil
	}
}

func RemoteShell(shell string) func(r *Rsync) error {
	return func(r *Rsync) error {
		if len(shell) == 0 {
			return errors.New("empty remote shell supplied")
		}

		s := fmt.Sprintf(`--rsh="%s"`, shell)
		r.args = append(r.args, s)
		return nil
	}
}
