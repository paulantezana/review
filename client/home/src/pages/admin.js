import React from "react";
import SiteLayout from "../layout/SiteLayout";
import { Router } from "@reach/router";
import PrivateRoute from "../components/PrivateRoute";
import Login from "../components/Login";

import Profile from "../components/admin/Profile";
import Admin from "../components/admin/Admin";
import Institution from "../components/admin/Institution";
import Module from "../components/admin/Module";
import { Layout } from "antd";

const AdminPage = () => (
    <SiteLayout>
        <Layout.Content className="Container" style={{padding: '24px 0'}}>
            <Router>
                <PrivateRoute path="/admin" component={Admin} />
                <PrivateRoute path="/admin/function" component={Module} />
                <PrivateRoute path="/admin/module" component={Module} />
                <PrivateRoute path="/admin/institution" component={Institution} />
                <PrivateRoute path="/admin/profile" component={Profile} />
                <Login path="/admin/login"/>
            </Router>
        </Layout.Content>
    </SiteLayout>
)
export default AdminPage