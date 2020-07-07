import React from 'react'
import Layout from '../components/layout'
import Head from 'next/head'

type Props = {
}

const Playing = () => {
  return (
    <React.Fragment>
      <Layout>
        <Head>
          <title>Street Fighter II</title>
        </Head>
        <h1>Playing Page</h1>
      </Layout>
    </React.Fragment>
  )
}

export default Playing

export const getStaticProps = async () => {
  return {
    props: {},
  }
}
