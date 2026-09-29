import GameList, { type GameListItem } from '../components/GameList'
import styles from './HomePage.module.css'

// Hand-picked games to feature on the homepage. Edit this array to change them.
const handPickedGames: GameListItem[] = [
  { appid: 292030, name: 'The Witcher 3: Wild Hunt - Complete Edition' },
  { appid: 1245620, name: 'ELDEN RING' },
  { appid: 553850, name: 'HELLDIVERS™ 2' },
  { appid: 322330, name: "Don't Starve Together" },
  { appid: 1623730, name: 'Palworld' },
  { appid: 1222670, name: 'The Sims™ 4' },
  { appid: 377160, name: 'Fallout 4' },
  { appid: 413150, name: 'Stardew Valley' },
  { appid: 108600, name: 'Project Zomboid' },
  { appid: 220, name: 'Half-Life 2' },
  { appid: 105600, name: 'Terraria' },
]

function HomePage() {
  return (
    <div className={styles.page}>
      <GameList title="Hand Picked" games={handPickedGames} />
    </div>
  )
}

export default HomePage
