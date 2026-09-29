import { Link } from 'react-router-dom'
import styles from './GameList.module.css'

export interface GameListItem {
  appid: number
  name: string
  icon?: string
}

interface GameListProps {
  title?: string
  games: GameListItem[]
}

function GameList({ title, games }: GameListProps) {
  return (
    <section className={styles.list}>
      {title && <h3 className={styles.title}>{title}</h3>}
      {games.length === 0 ? (
        <p className={styles.empty}>No games to show yet.</p>
      ) : (
        <ul className={styles.items}>
          {games.map((game) => (
            <li key={game.appid}>
              <Link to={`/game/${game.appid}`} className={styles.item}>
                {game.icon ? (
                  <img src={game.icon} alt="" className={styles.icon} />
                ) : (
                  <span className={styles.iconPlaceholder} aria-hidden="true" />
                )}
                <span className={styles.name}>{game.name}</span>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}

export default GameList
