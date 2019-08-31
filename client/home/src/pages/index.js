import React from "react"
import { graphql, Link } from "gatsby"
import SiteLayout from "../layout/SiteLayout"
import MainBanner from "../components/MainBanner"
import Instances  from '../components/Instances';

import ImageHome from '../images/app-responsive.png';

import { Button, List } from "antd"
import SEO from "../components/seo"

import { Modules } from '../data/modules';

export default ({ data }) => (
    <SiteLayout>
        <SEO title="Introduccion" />

        <MainBanner/>

        <Instances/>

        <div className="BannerHome">
            <div className="Container Grid m-2">
                <div className="BannerHome-left">
                    <h2 className="BannerHome-title">
                        Sistema Institucional
                        <div>Web</div>
                    </h2>
                    <p>lorem</p>
                    <Button type="primary">Documentacion</Button>
                </div>
                <div className="BannerHome-right">
                    <img src={ImageHome} alt="banner" />
                </div>
            </div>
        </div>

        <div className="Container SnMt64 SnMb64">
            <List
                grid={{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 4, xl: 4, xxl: 4 }}
                dataSource={Modules}
                renderItem={item => (
                    <List.Item>
                        <div>
                            <h3>{ item.title }</h3>
                            <p>{ item.text }</p>
                            <Link to={item.doc}>Seguir leendo</Link>
                        </div>
                    </List.Item>
                )}
            />
        </div>
    </SiteLayout>
)

export const query = graphql`
    query {
        site {
            siteMetadata {
                title
            }
        }
    }
`
