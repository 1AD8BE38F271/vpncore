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

package conn

import (
	"net"
)

func Dial(contexts []ConnLayerContext) (c net.Conn, err error) {
	if len(contexts) < 1 {
		return nil, ErrInvalidArgs
	}

	ctx := contexts[0]
	if ctx.Layer() != STREAM_LAYER {
		return nil, ErrInvalidCtx
	}

	c, err = ctx.Dial(nil)
	if err != nil {
		return
	}

	for _, ctx := range contexts[1:] {
		c, err = ctx.Dial(c)
		if err != nil {
			return
		}

	}

	return c, err
}

