import React, { Component } from 'react';
import './Board.scss';

class Item extends Component {
  render() {
    const { row, col, val } = this.props;
    if (val === '.') {
      return <span className="crossword-board__item--blank" id={`item-${row}-${col}`} />;
    }
    return (
      <input
        id={`item-${row}-${col}`}
        className="crossword-board__item"
        type="text"
        minLength="1"
        maxLength="1"
        required="required"
        value=""
        readOnly
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
        const idx = row * game.grid_rows + (col % game.grid_cols);
        items.push(<Item key={idx} row={row + 1} col={col + 1} val={game.grid[idx]} />);
      }
    }

    let labels = [];
    for (let row = 0; row < game.grid_rows; row++) {
      for (let col = 0; col < game.grid_cols; col++) {
        const idx = row * game.grid_rows + (col % game.grid_cols);
        const label = game.grid_nums[idx];
        if (label === 0) continue;
        labels.push(
          <span
            key={`label-${label}`}
            style={{
              gridColumnStart: col + 1,
              gridColumnEnd: col + 1,
              gridRowStart: row + 1,
              gridRowEnd: row + 1
            }}
            className={`crossword-board__item-label`}
          >
            <span className="crossword-board__item-label-text">{label}</span>
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
        <div className="crossword-board-container">
          <div className="crossword-board" style={crosswordBoardStyle}>
            {items}
            <div className="crossword-board crossword-board--highlight" style={crosswordBoardStyle}>
              <span
                className="crossword-board__item-highlight"
                style={{
                  gridColumnStart: 1,
                  gridColumnEnd: 8,
                  gridRowStart: 1,
                  gridRowEnd: 1
                }}
              />
            </div>

            <div className="crossword-board crossword-board--labels" style={crosswordBoardStyle}>
              {labels}
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default Board;
