import Link from 'next/link'
import styles from './game-card.module.css'

const JoinableGameCard = () => {
  return (
    <div className={styles['card-frame']}>
      <div className="relative">
        <div className={styles['btn-container']}>
          <Link href="/playing">
            <a className={styles['btn']}>Join Game!</a>
          </Link>
        </div>      

        <img className="cover-image w-full h-28 object-center object-cover"
          src="https://tailwindcss.com/img/tailwind-ui-sidebar.png" />

        <div className="p-2">
          <div className="text-sm">
            <span className="font-bold">Game Title</span>
          </div>
          <div>
            <span className="text-sm">1</span>
            <span className="text-xs mx-1">/</span>
            <span className="text-sm">2</span>
          </div>
        </div>
      </div>
    </div>
  )
}

export default JoinableGameCard
