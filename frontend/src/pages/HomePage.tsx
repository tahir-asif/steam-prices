import GameList, { type GameListItem } from '../components/GameList'
import styles from './HomePage.module.css'

interface GameListSection {
  title: string
  games: GameListItem[]
}

const sections: GameListSection[] = [
  {
    title: 'Most Popular',
    games: [
      { appid: 730, name: 'Counter-Strike 2' },
      { appid: 570, name: 'Dota 2' },
      { appid: 578080, name: 'PUBG: Battlegrounds' },
      { appid: 271590, name: 'Grand Theft Auto V' },
      { appid: 1172470, name: 'Apex Legends' },
      { appid: 252490, name: 'Rust' },
      { appid: 1085660, name: 'Destiny 2' },
      { appid: 1245620, name: 'ELDEN RING' },
    ],
  },
  {
    title: 'Cheapest',
    games: [
      { appid: 620, name: 'Portal 2' },
      { appid: 105600, name: 'Terraria' },
      { appid: 413150, name: 'Stardew Valley' },
      { appid: 367520, name: 'Hollow Knight' },
      { appid: 504230, name: 'Celeste' },
      { appid: 1145360, name: 'Hades' },
      { appid: 268910, name: 'Cuphead' },
      { appid: 391540, name: 'Undertale' },
    ],
  },
  {
    title: 'Biggest Discounts',
    games: [
      { appid: 1091500, name: 'Cyberpunk 2077' },
      { appid: 1174180, name: 'Red Dead Redemption 2' },
      { appid: 359550, name: "Tom Clancy's Rainbow Six Siege" },
      { appid: 552520, name: 'Far Cry 5' },
      { appid: 812140, name: "Assassin's Creed Odyssey" },
      { appid: 377160, name: 'Fallout 4' },
      { appid: 489830, name: 'The Elder Scrolls V: Skyrim' },
      { appid: 1238840, name: 'Battlefield 1' },
    ],
  },
  {
    title: 'New Releases',
    games: [
      { appid: 2358720, name: 'Black Myth: Wukong' },
      { appid: 1086940, name: "Baldur's Gate 3" },
      { appid: 1623730, name: 'Palworld' },
      { appid: 1966720, name: 'Lethal Company' },
      { appid: 553850, name: 'HELLDIVERS 2' },
      { appid: 1778820, name: 'TEKKEN 8' },
      { appid: 1904540, name: 'Enshrouded' },
      { appid: 1364780, name: 'Street Fighter 6' },
    ],
  },
]

function HomePage() {
  return (
    <div className={styles.page}>
      <div className={styles.lists}>
        {sections.map((section) => (
          <GameList
            key={section.title}
            title={section.title}
            games={section.games}
          />
        ))}
      </div>
    </div>
  )
}

export default HomePage
