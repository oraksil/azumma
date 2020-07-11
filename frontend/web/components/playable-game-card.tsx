import Link from 'next/link'
import styles from './game-card.module.css'

const PlayableGameCard = () => {
  return (
    <div className={styles['card-frame']}>
      <img className="cover-image w-full h-40 object-center object-cover"
        src="https://tailwindcss.com/img/tailwind-ui-sidebar.png" />

      <div className="relative">
        <div className={styles['btn-container']}>
          <Link href="/playing">
            <a className={styles['btn']}>Play Game!</a>
          </Link>
        </div>
        <div className="p-2">
          <div className="text-sm">
            <span className="font-bold">Game Title</span>
          </div>
          <div className="text-xs">
             <span>Game Producer</span>
          </div>
          <div>
            <span className="text-sm">2</span>
            <span className="text-xs">P</span>
          </div>
        </div>
      </div>
    </div>
  )
}

export default PlayableGameCard 
