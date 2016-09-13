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
 */

package cmd

import (
	"strings"
	"os/exec"
	"fmt"
)

func RunCommand(cmd string) (out string, err error) {

	args := strings.Split(cmd, " ")
	command := exec.Command(args[0], args[1:]...)

	out_bytes, err := command.CombinedOutput()
	out = string(out_bytes)
	if err != nil {
		fmt.Printf("RUN[%s]==>\n %s==================\n", cmd, out)
		return
	}
	fmt.Printf("RUN[%s]==>\n %s==================\n", cmd, out)
	return
}
