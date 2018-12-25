import React, { Component } from 'react';
import './Board.scss';

class Item extends Component {
  render() {
    const { row, col, val } = this.props;
    if (val === '.') {
      return <span class="crossword-board__item--blank" id={`item-${row}-${col}`} />;
    }
    return (
      <input
        id={`item-${row}-${col}`}
        class="crossword-board__item"
        type="text"
        minlength="1"
        maxlength="1"
        required="required"
        value=""
      />
    );
  }
}

class Board extends Component {
  render() {
    const { game } = this.props;
    let items = [];
    for (let row = 0; row < game.grid_rows; row++) {
      for (let col = 0; col < game.grid_cols; col++) {
        items.push(<Item row={row + 1} col={col + 1} val={game.grid[row * game.grid_rows + (col % game.grid_cols)]} />);
      }
    }

    let labels = [];
    for (let row = 0; row < game.grid_rows; row++) {
      for (let col = 0; col < game.grid_cols; col++) {
        const label = game.grid_nums[row * game.grid_rows + (col % game.grid_cols)];
        console.log(label);
        if (label === 0) continue;
        labels.push(
          <span
            id={`label-${label}`}
            style={{
              gridColumnStart: col + 1,
              gridColumnEnd: col + 1,
              gridRowStart: row + 1,
              gridRowEnd: row + 1
            }}
            class={`crossword-board__item-label`}
          >
            <span class="crossword-board__item-label-text">{label}</span>
          </span>
        );
      }
    }

    const crosswordBoardStyle = {
      width: 50 * game.grid_cols + 'px',
      height: 50 * game.grid_rows + 'px',
      gridTemplateRows: 'repeat(' + game.grid_rows + ', ' + 100 / game.grid_rows + '%)',
      gridTemplateColumns: 'repeat(' + game.grid_cols + ', ' + 100 / game.grid_cols + '%)'
    };

    return (
      <div>
        <div class="crossword-board-container">
          <div class="crossword-board" style={crosswordBoardStyle}>
            {items}
            <div class="crossword-board crossword-board--highlight" style={crosswordBoardStyle}>
              <span
                class="crossword-board__item-highlight"
                style={{
                  gridColumnStart: 1,
                  gridColumnEnd: 8,
                  gridRowStart: 1,
                  gridRowEnd: 1
                }}
              />
            </div>

            <div class="crossword-board crossword-board--labels" style={crosswordBoardStyle}>
              {labels}
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default Board;
