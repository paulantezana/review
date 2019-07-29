import React from "react"
import SiteLayout from "../layout/SiteLayout"
import { Button } from "antd"
import QueueAnim from "rc-queue-anim"
import BannerImage from "../components/BannerImage"
import { Link } from "gatsby"
import SEO from "../components/seo"

export default ({ data }) => (
    <SiteLayout>
        <SEO
            title="Documentacion Codigo Fuente"
            description="Documentación de código fuente del sistema cualquier cambio que realice en el sistema debe documentar en la API"
        />
        <div className="ApiBanner Container">
            <QueueAnim>
                <h1 key="h2"> Documentacion Codigo Fuente</h1>
                <p key="p">
                    Documentación de código fuente del sistema cualquier cambio
                    que realice en el sistema debe documentar en la API
                </p>
                <span key="button">
                    <Link to="/documentacion/starter">
                        <Button type="primary">Leer documentacion</Button>
                    </Link>
                </span>
            </QueueAnim>
            <BannerImage />
        </div>
    </SiteLayout>
)
