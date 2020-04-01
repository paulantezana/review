import React from "react";
import SiteLayout from "../layout/SiteLayout";
import { Button } from "antd";
import QueueAnim from "rc-queue-anim";
import { Link } from "gatsby";
import SEO from "../components/seo";
import slideHome from "../images/slide-home.jpg";

export default ({ data }) => (
    <SiteLayout>
        <SEO
            title="Tutoriales"
            description="Documentación de código fuente del sistema cualquier cambio que realice en el sistema debe documentar en la API"
        />
        <div className="DocBanner">
            <img src={slideHome} className="DocBanner-bg" alt="slide-bg"/>
            <div className="Container">
                <QueueAnim className="DocBanner-data">
                    <h1> Tutoriales</h1>
                    <p>
                        Cursos disponibles en donde aprenderás a controlar el sistema y
                        personalizas de acuerdo a tus necesidades. 
                    </p>
                    <Link to="/tutorial/starter">
                        <Button>Ver tutoriales</Button>
                    </Link>
                </QueueAnim>
                <div>

                </div>
            </div>
        </div>
    </SiteLayout>
)
