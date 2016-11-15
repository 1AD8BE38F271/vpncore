/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Author: FTwOoO <booobooob@gmail.com>
 * Modification Author: 1AD8BE38F271 <1AD8BE38F271@protonmail.com>
 */

package cmd

import (
	"strings"
	"os/exec"

    . "github.com/1AD8BE38F271/vpncore"
)

func RunCommand(cmd string) (out string, err error) {
	args := strings.Split(cmd, " ")
	command := exec.Command(args[0], args[1:]...)

	out_bytes, err := command.CombinedOutput()
	out = string(out_bytes)
	if err != nil {
		Logger.Warningf("RUN[%s] ==> %s", cmd, out)
		return
	}
	Logger.Infof("RUN[%s] ==> %s", cmd, out)
	return
}
