import React, { Component } from "react";
import { navigate } from "gatsby";
import { getToken } from "../utils/authority";

class PrivateRoute extends Component {
    componentDidMount() {
        const { location } = this.props
        let noOnLoginPage = location.pathname !== `/admin/login`
        if (!getToken() && noOnLoginPage) {
            navigate("/admin/login")
            return null
        }
    }
    render() {
        const { component: Component, ...rest } = this.props
        return <Component {...rest} />
    }
}

export default PrivateRoute