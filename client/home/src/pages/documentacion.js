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
        <div className="DocBanner">
            <div className="Container">
                <QueueAnim className="DocBanner-data">
                    <h1> Documentacion Codigo Fuente</h1>
                    <p>
                        Documentación de código fuente del sistema cualquier cambio
                        que realice en el sistema debe documentar en la API
                    </p>
                    <Link to="/documentacion/starter">
                        <Button type="primary">Leer documentacion</Button>
                    </Link>
                </QueueAnim>
                <BannerImage />
            </div>
        </div>
    </SiteLayout>
)
