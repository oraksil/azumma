import Head from 'next/head'
import Link from 'next/link'
import Layout from '../components/layout'
import PlayableGameCard from '../components/playable-game-card'

type Props = {
}

const Catalog  = () => {
  const games = [1, 2, 3]

  return (
    <Layout>
      <Head>
        <title>Go enjoy with games!</title>
      </Head>
      <div className="container pt-32">
        <div className="game-cards flex justify-center">
          {games.map((_, i) => ( 
            <div key={i} className="game-card inline-block mx-4">
              <PlayableGameCard />
            </div>
          ))}
        </div>
        <div className="existing-games flex justify-center pt-4">
          <Link href="/hall">
            <a>
              <span className="text-sm underline">Join existing games...</span>
            </a>
          </Link>
        </div>
      </div>
    </Layout>
  )
}

export default Catalog

export const getStaticProps = async () => {
  return {
    props: {},
  }
}
