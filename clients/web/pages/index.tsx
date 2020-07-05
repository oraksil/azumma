import React from 'react'
import Container from '../components/container'
import Layout from '../components/layout'
import Head from 'next/head'

type Props = {
}

const Index = () => {
  return (
    <React.Fragment>
      <Layout>
        <Head>
          <title>Welcome to Oraksil!</title>
        </Head>
        <Container>
          <h1>Landing Page</h1>
        </Container>
      </Layout>
    </React.Fragment>
  )
}

export default Index

export const getStaticProps = async () => {
  return {
    props: {},
  }
}
