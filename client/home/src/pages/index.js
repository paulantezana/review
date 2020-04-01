import React from "react"
import { graphql } from "gatsby"
import SiteLayout from "../layout/SiteLayout"
import MainBanner from "../components/Home/MainBanner"
import Instances  from '../components/Home/Instances';

import ImageHome from '../images/app-responsive.png';

import { Button } from "antd"
import SEO from "../components/seo"

// import { Modules } from '../data/modules';
import Modules from '../components/Home/Modules';

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
            <Modules/>
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
