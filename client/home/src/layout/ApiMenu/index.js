import React from "react"
import { Menu, Icon, Affix } from "antd"
import { Link } from "gatsby"
import {
    IconAdmission,
    IconReview,
    IconCertificate,
    IconLibrarie,
    IconMonitoring,
    IconWebSite,
} from "../../components/icons"

const SubMenu = Menu.SubMenu

class Sider extends React.Component {
    render() {
        return (
            <div>
                <Affix offsetTop={0}>
                    <Menu mode="inline" className="SiderMenu">
                        <Menu.Item key="starter">
                            <Icon type="bulb" />
                            <span>
                                <Link to="/documentacion/starter">
                                    Introducción
                                </Link>
                            </span>
                        </Menu.Item>
                        <Menu.Item key="prerequisitos">
                            <Icon type="sync" />
                            <span>
                                <Link to="/documentacion/prerequisitos">
                                    Prerequisitos
                                </Link>
                            </span>
                        </Menu.Item>
                        <SubMenu
                            key="instalar"
                            title={
                                <span>
                                    <Icon type="code" />
                                    <span>Instalar</span>
                                </span>
                            }
                        >
                            <Menu.Item>
                                <Link to="/documentacion/instalar/cliente">
                                    Cliente
                                </Link>
                            </Menu.Item>
                            <Menu.Item>
                                <Link to="/documentacion/instalar/servidor">
                                    Servidor
                                </Link>
                            </Menu.Item>
                        </SubMenu>
                        <SubMenu
                            key="estrutura"
                            title={
                                <span>
                                    <Icon type="folder" />
                                    <span>Estrutura</span>
                                </span>
                            }
                        >
                            <Menu.Item>
                                <Link to="/documentacion/estrutura/cliente">
                                    Cliente
                                </Link>
                            </Menu.Item>
                            <Menu.Item>
                                <Link to="/documentacion/estrutura/servidor">
                                    Servidor
                                </Link>
                            </Menu.Item>
                        </SubMenu>
                        {/* <Menu.Item key="arquitectura">Arquitectura</Menu.Item>
                        <Menu.Item key="basededatos">Base de datos</Menu.Item> */}
                        <SubMenu
                            key="admision"
                            title={
                                <span>
                                    <Icon
                                        component={IconAdmission}
                                        style={{
                                            fontSize: "21px",
                                            position: "relative",
                                            top: "5px",
                                        }}
                                    />
                                    <span>Admisión</span>
                                </span>
                            }
                        >
                            <Menu.Item>
                                <Link to="/documentacion/admision/instituto">
                                    Instituto
                                </Link>
                            </Menu.Item>
                            <Menu.Item>
                                <Link to="/documentacion/admision/modelos">
                                    Modelos
                                </Link>
                            </Menu.Item>
                            <Menu.Item>
                                <Link to="/documentacion/admision/controladores">
                                    Controladores
                                </Link>
                            </Menu.Item>
                            <Menu.Item>
                                <Link to="/documentacion/admision/cliente">
                                    Cliente
                                </Link>
                            </Menu.Item>
                            <Menu.Item>
                                <Link to="/documentacion/admision/softwares">
                                    Softwares
                                </Link>
                            </Menu.Item>
                        </SubMenu>
                        <SubMenu
                            key="revision"
                            title={
                                <span>
                                    <Icon
                                        component={IconReview}
                                        style={{
                                            fontSize: "21px",
                                            position: "relative",
                                            top: "5px",
                                        }}
                                    />{" "}
                                    <span>Revisión</span>
                                </span>
                            }
                        >
                            <Menu.Item>
                                <Link to="/documentacion/revision/modelos">
                                    Modelos
                                </Link>
                            </Menu.Item>
                            <Menu.Item>
                                <Link to="/documentacion/revision/softwares">
                                    Softwares
                                </Link>
                            </Menu.Item>
                        </SubMenu>
                        <SubMenu
                            key="certificacion"
                            title={
                                <span>
                                    <Icon
                                        component={IconCertificate}
                                        style={{
                                            fontSize: "21px",
                                            position: "relative",
                                            top: "5px",
                                        }}
                                    />{" "}
                                    <span>Certificación</span>
                                </span>
                            }
                        >
                            <Menu.Item>
                                <Link to="/documentacion/certificacion/modelos">
                                    Modelos
                                </Link>
                            </Menu.Item>
                            <Menu.Item>
                                <Link to="/documentacion/certificacion/softwares">
                                    Softwares
                                </Link>
                            </Menu.Item>
                        </SubMenu>
                        <SubMenu
                            key="egresados"
                            title={
                                <span>
                                    <Icon
                                        component={IconMonitoring}
                                        style={{
                                            fontSize: "21px",
                                            position: "relative",
                                            top: "5px",
                                        }}
                                    />{" "}
                                    <span>Egresados</span>
                                </span>
                            }
                        >
                            <Menu.Item>
                                <Link to="/documentacion/egresados/modelos">
                                    Modelos
                                </Link>
                            </Menu.Item>
                            <Menu.Item>
                                <Link to="/documentacion/egresados/softwares">
                                    Softwares
                                </Link>
                            </Menu.Item>
                        </SubMenu>
                        <SubMenu
                            key="biblioteca"
                            title={
                                <span>
                                    <Icon
                                        component={IconLibrarie}
                                        style={{
                                            fontSize: "21px",
                                            position: "relative",
                                            top: "5px",
                                        }}
                                    />{" "}
                                    <span>Biblioteca</span>
                                </span>
                            }
                        >
                            <Menu.Item>
                                <Link to="/documentacion/biblioteca/modelos">
                                    Modelos
                                </Link>
                            </Menu.Item>
                            <Menu.Item>
                                <Link to="/documentacion/biblioteca/softwares">
                                    Softwares
                                </Link>
                            </Menu.Item>
                        </SubMenu>

                        <SubMenu
                            key="sitioweb"
                            title={
                                <span>
                                    <Icon
                                        component={IconWebSite}
                                        style={{
                                            fontSize: "21px",
                                            position: "relative",
                                            top: "5px",
                                        }}
                                    />{" "}
                                    <span>Sitio web</span>
                                </span>
                            }
                        >
                            <Menu.Item>
                                <Link to="/documentacion/sitioweb/tema">
                                    Theme
                                </Link>
                            </Menu.Item>
                            <Menu.Item>
                                <Link to="/documentacion/sitioweb/plugin">
                                    Plugin
                                </Link>
                            </Menu.Item>
                            <Menu.Item>
                                <Link to="/documentacion/sitioweb/softwares">
                                    Softwares
                                </Link>
                            </Menu.Item>
                        </SubMenu>
                        <SubMenu
                            key="mensajeria"
                            title={
                                <span>
                                    <Icon type="rest" />
                                    <span>Mensajería</span>
                                </span>
                            }
                        >
                            <Menu.Item>
                                <Link to="/documentacion/mensajeria/modelos">
                                    Modelos
                                </Link>
                            </Menu.Item>
                            <Menu.Item>
                                <Link to="/documentacion/mensajeria/softwares">
                                    Softwares
                                </Link>
                            </Menu.Item>
                        </SubMenu>
                        <SubMenu
                            key="alumno"
                            title={
                                <span>
                                    <Icon type="rest" />
                                    <span>Alumno</span>
                                </span>
                            }
                        >
                            <Menu.Item>
                                <Link to="/documentacion/alumno/modelos">
                                    Modelos
                                </Link>
                            </Menu.Item>
                            <Menu.Item>
                                <Link to="/documentacion/alumno/softwares">
                                    Softwares
                                </Link>
                            </Menu.Item>
                        </SubMenu>
                    </Menu>
                </Affix>
            </div>
        )
    }
}

export default Sider
