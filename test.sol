pragma solidity ^0.7.0;
import "./Table.sol";

contract supply_chain{
    
    struct Bill {
        string borrer;
        int amount;
    }

    struct company {
        string name;               //名称
        address addr;              //地址
        int asset;                 //资产
        Bill[] bills;              //作为放款人的账单
    }








    event Pay_query(address addr, uint amount);
    event Receive_query(address addr, uint amount);
    event New_company(string name, address addr);
    event New_bank(address addr);
    event issue(address to, uint account); 

    //根据id获取应还账款
    function getPay(uint id)  public view returns(uint, address, uint, uint) {
        Company storage company = companies[msg.sender];
        Pay storage pay = company.pays[id];
        return (pay.id, pay.creditor, pay.amount, pay.ddl;
    }

    //根据id获取应收账款
    function getReceive(uint id) public view returns(uint, address, uint, uint) {
        Company storage company = companies[msg.sender];
        Receive storage receive = company.receives[id];
        return (receive.id, receive.debtor, receive.amount, receive.ddl);
    }

    //获取一个公司总的应还账款（总负债额）
    function get_total_pay(address addr) public view return(uint) {
        uint total = 0;
        for(int i = 0; i < companies[addr].pays.length; ++i)
        {
            total += companies[addr].pays[i];
        }
        emit Pay_query(addr, total);
        return total;
    }

    //获取一个公司总的应收账款（总放债额）
    function get_total_receive(address addr) public view return(uint) {
        uint total = 0;
        for(int i = 0; i < companies[addr].receives.length; ++i)
        {
            total += companies[addr].receives[i];
        }
        emit Recieve_query(addr, total);
        return total;
    }

    //注册公司
    function register_company(string name) public {
        Company storage company = companies[msg.sender];
        company.name = name;
        company.addr = msg.sender;
        emit New_company(name, msg.sender);
    }

    //添加银行
    function add_bank(address addr) public {
        for (uint i = 0; i < banks.length; ++i) {
            if (msg.sender == banks[i]) {
                banks.push(addr);
                emit New_bank(addr);
            }
        }
    }

    //判断是否为银行
    function isBank() public view returns (bool) {
        for (uint i = 0; i < banks.length; ++i) {
            if (msg.sender == banks[i]) {
                return true;
            }
        }
        return false;
    }

    //向指定方发行货币
    function issue(address to, uint amount) public {
        for (uint i = 0; i < banks.length; ++i) {
            if (msg.sender == banks[i]) {
                companies[to].assets += amount;
                emit Issue(to, amount);
            }
        }
    }



}