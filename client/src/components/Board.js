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
        value={val}
        readOnly
      />
    );
  }
}

class Board extends Component {
  render() {
    const { game } = this.props;

    let items = [];
    for (let row = 0; row < game.layout.size.rows; row++) {
      for (let col = 0; col < game.layout.size.cols; col++) {
        const idx = row * game.layout.size.rows + (col % game.layout.size.cols);
        items.push(<Item key={idx} row={row + 1} col={col + 1} val={game.layout.grid[idx]} />);
      }
    }

    let labels = [];
    for (let row = 0; row < game.layout.size.rows; row++) {
      for (let col = 0; col < game.layout.size.cols; col++) {
        const idx = row * game.layout.size.rows + (col % game.layout.size.cols);
        const label = game.layout.gridnums[idx];
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
      width: 50 * game.layout.size.cols + 'px',
      height: 50 * game.layout.size.rows + 'px',
      gridTemplateRows: 'repeat(' + game.layout.size.rows + ', ' + 100 / game.layout.size.rows + '%)',
      gridTemplateColumns: 'repeat(' + game.layout.size.cols + ', ' + 100 / game.layout.size.cols + '%)'
    };

    let highlight_style = { opacity: 0 };
    if (game.current_clue) {
      const current_clue_label = game.current_clue.description.split('.')[0];
      const current_clue_direction = game.current_clue.direction;
      let current_clue_row = 1;
      let current_clue_col = 1;
      for (let row = 0; row < game.layout.size.rows; row++) {
        for (let col = 0; col < game.layout.size.cols; col++) {
          const idx = row * game.layout.size.rows + (col % game.layout.size.cols);
          if (game.layout.gridnums[idx] === parseInt(current_clue_label)) {
            current_clue_row = row + 1;
            current_clue_col = col + 1;
          }
        }
      }
      if (current_clue_direction === 'across') {
        highlight_style = {
          gridColumnStart: current_clue_col,
          gridColumnEnd: current_clue_col + game.current_clue.answer.length,
          gridRowStart: current_clue_row,
          gridRowEnd: current_clue_row
        };
      } else {
        highlight_style = {
          gridColumnStart: current_clue_col,
          gridColumnEnd: current_clue_col,
          gridRowStart: current_clue_row,
          gridRowEnd: current_clue_row + game.current_clue.answer.length
        };
      }
    }

    return (
      <div>
        <div className="crossword-board-container">
          <div className="crossword-board" style={crosswordBoardStyle}>
            {items}
            <div className="crossword-board crossword-board--highlight" style={crosswordBoardStyle}>
              <span className="crossword-board__item-highlight" style={highlight_style} />
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
