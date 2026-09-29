import { useState } from 'react'
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

const steamCapsule = (appid: number) =>
  `https://cdn.cloudflare.steamstatic.com/steam/apps/${appid}/capsule_231x87.jpg`

function GameThumb({ appid, icon }: { appid: number; icon?: string }) {
  const [failed, setFailed] = useState(false)

  if (failed) {
    return <span className={styles.iconPlaceholder} aria-hidden="true" />
  }

  return (
    <img
      src={icon ?? steamCapsule(appid)}
      alt=""
      loading="lazy"
      className={styles.icon}
      onError={() => setFailed(true)}
    />
  )
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
                <GameThumb appid={game.appid} icon={game.icon} />
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
