import React from 'react'
import Container from '../components/container'
import Layout from '../components/layout'
import Head from 'next/head'

type Props = {
}

const Catalog  = () => {
  return (
    <React.Fragment>
      <Layout>
        <Head>
          <title>Go enjoy with games!</title>
        </Head>
        <Container>
          <h1>Catalog Page</h1>
        </Container>
      </Layout>
    </React.Fragment>
  )
}

export default Catalog

export const getStaticProps = async () => {
  return {
    props: {},
  }
}
