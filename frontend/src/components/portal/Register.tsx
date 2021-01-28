import * as React from "react";
import AppContext from "../../AppContext";
import { DefaultButton, PrimaryButton, Stack, Text, TextField } from "@fluentui/react";
import useHttp from "../../hooks/http";
import UserInfo from "../../models/UserInfo";

export interface RegisterModel {
  setType: (type: string) => void;
}

const Register: React.FunctionComponent<RegisterModel> = (props) => {
  const { setUser } = React.useContext(AppContext);
  const [name, setName] = React.useState<string>();
  const [password, setPassword] = React.useState<string>();
  const [email, setEmail] = React.useState<string>();
  const [error, setError] = React.useState<string>();
  const registerRequest = useHttp<{ status: string; }>("/api/user/register", "POST");
  const userInfoRequest = useHttp<UserInfo>("/api/user/self", "GET");

  React.useEffect(() => {
    if (userInfoRequest.data && !userInfoRequest.loading && setUser) {
      setUser({
        address: userInfoRequest.data.address,
        username: userInfoRequest.data.username,
      });
    }
  }, [userInfoRequest.loading, userInfoRequest.data]);

  React.useEffect(() => {
    if (!registerRequest.loading) {
      if (registerRequest.data?.status === "success") {
        console.log("注册成功");
        userInfoRequest.fire();
        setError(undefined);
      }
      else if (registerRequest.data || !registerRequest.ok) {
        setError("注册失败");
      }
    }
  }, [registerRequest.loading, registerRequest.data, registerRequest.ok]);

  const register = () => {
    registerRequest.fire({
      username: name,
      password: password,
      email: email
    });
  };

  return <Stack>
    <Stack.Item>
      <Text variant="xxLarge">注册账户</Text>
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <TextField label="用户名" defaultValue={name} onChange={(_, v) => setName(v)} />
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <TextField label="密码" canRevealPassword={true} type="password" defaultValue={password} onChange={(_, v) => setPassword(v)} />
    </Stack.Item>
    {/* <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <TextField label="邮箱" type="email" defaultValue={email} onChange={(_, v) => setEmail(v)} />
    </Stack.Item> */}
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <Stack horizontal>
        <Stack.Item>
          <PrimaryButton text="注册" onClick={(register)} />
        </Stack.Item>
        <Stack.Item styles={{ root: { paddingLeft: 10 } }}>
          <DefaultButton text="登录" onClick={() => props.setType("login")} />
        </Stack.Item>
      </Stack>
    </Stack.Item>
    {
      error && <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
        <p style={{ color: "red" }}>{error}</p>
      </Stack.Item>
    }
  </Stack >;
};

export default Register;