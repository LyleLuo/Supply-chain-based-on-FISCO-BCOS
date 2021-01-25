pragma solidity ^0.4.25;

contract supply_chain {
    //欠条
    struct Debt {
        uint id;
        uint credior_id;   //此账单在对应债务方的收据中的索引
        address creditor;
        uint amount;
    }
    
    //收据
    struct Receipt {
        uint id;
        uint borrer_id;    //此账单在对应欠债方的欠条中的索引    
        address borrer;
        uint amount;
    }

    struct Company {
        string name;               //名称
        address addr;              //地址
        uint asset;                //资产
        Debt[] debts;              //欠条
        Receipt[] receipts;        //收据
    }

    struct Bank {
        string name;
        address addr;
    }

    event Debt_query(string name, uint amount);
    event Receipt_query(string name, uint amount);
    event Register_company(string name, address addr, uint asset);
    event Issue(address to, uint account); 
    event Create_bill(address from, address to, uint amount);
    event Transfer_bill(address from, address to, uint from_id, uint to_id, uint amount);
    event Finance(address to, uint amount);
    event Settle(address to, uint bill_id);
    
    mapping(address => Company) public companies;
    Bank public bank;
    
    //构造函数初始化唯一的中央银行，有权发行信用
    constructor() {
        bank.name = "CentralBank";
        bank.addr = msg.sender;
    }

    function bank_info() public returns(string memory bank_name,address bank_address) {
        return (bank.name, bank.addr);
    }

    //注册公司
    function register_company(string memory name, uint asset) public {
        Company storage company = companies[msg.sender];
        company.name = name;
        company.addr = msg.sender;
        company.asset = asset;
        emit Register_company(name, msg.sender, asset);
    }

    //查询公司总资产
    function getAsset() public returns(uint asset) {
        return companies[msg.sender].asset;
    }

    //根据账单获取应还账款
    function getDebt(uint bill_id) public returns(address creditor, uint amount) {
        require(bill_id < companies[msg.sender].debts.length, "the query bill does not exist");
        address tem_creditor = companies[msg.sender].debts[bill_id].creditor;
        uint tem_amount = companies[msg.sender].debts[bill_id].amount;
        return (tem_creditor, tem_amount);
    }

    //根据账单获取应收账款
    function getReceipt(uint bill_id) public returns(address borrer, uint amount) {
        require(bill_id < companies[msg.sender].receipts.length, "the query bill does not exist");
        address tem_borrer = companies[msg.sender].receipts[bill_id].borrer;
        uint tem_amount = companies[msg.sender].receipts[bill_id].amount;
        return (tem_borrer, tem_amount);
    }
    
    //获取一个公司总的应还账款（总负债额）
    function get_total_debt() public returns(uint total_debt) {
        uint total = 0;
        for(uint i = 0; i < companies[msg.sender].debts.length; ++i)
        {
            total += companies[msg.sender].debts[i].amount;
        }
        emit Debt_query(companies[msg.sender].name, total);
        return total;
    }    

    //获取一个公司总的应收账款（总放债额）
    function get_total_receipt() public returns(uint total_receipt) {
        uint total = 0;
        for(uint i = 0; i < companies[msg.sender].receipts.length; ++i)
        {
            total += companies[msg.sender].receipts[i].amount;
        }
        emit Receipt_query(companies[msg.sender].name, total);
        return total;
    }

    //向指定方发行货币
    function issue(address to, uint amount) public {
        require(msg.sender == bank.addr, "only bank could issue money");
        companies[to].asset += amount;
        emit Issue(to, amount);
    }
    
    //签发账单
    function create_bill(address from, address to, uint amount) public {
        require(msg.sender == bank.addr, "only bank could create bill");
        uint borrer_id = companies[to].debts.length;
        uint credior_id = companies[from].receipts.length;
        
        companies[from].receipts.push(
            Receipt(credior_id, borrer_id, to, amount)
        );
        companies[to].debts.push(
            Debt(borrer_id, credior_id, from, amount)  
        );
        emit Create_bill(from, to, amount);
    }
    
    //应收账款转移，涉及到三个人和两个账单，from_id是钱款来源账单(receipts中)，to_id是钱款去向账单（debts中）
    function transfer_bill(uint from_id, uint to_id, address to, uint amount) public {
        address before_addr = companies[msg.sender].receipts[from_id].borrer;
        //from_id <=> borrer_id, to_id <=> credior_id
        uint borrer_id = companies[msg.sender].receipts[from_id].borrer_id;
        uint credior_id = companies[msg.sender].debts[to_id].credior_id;
        
        //若A欠B50，B欠C100，转移50后，A欠C50, B欠C50，A不再欠B
        if(amount <= companies[msg.sender].receipts[from_id].amount) {
            companies[before_addr].debts[borrer_id].amount -= amount;
            if(companies[before_addr].debts[borrer_id].amount == 0)
            {
                delete companies[before_addr].debts[borrer_id];
            }
            
            companies[msg.sender].receipts[from_id].amount -= amount;
            if(companies[msg.sender].receipts[from_id].amount == 0)
            {
                delete companies[msg.sender].receipts[from_id];
            }
            
            companies[msg.sender].debts[to_id].amount -= amount;
            if(companies[msg.sender].debts[to_id].amount == 0) {
                delete companies[msg.sender].debts[to_id];
            }
            
            companies[to].receipts[credior_id].amount -= amount;
            if(companies[to].receipts[credior_id].amount == 0) {
                delete companies[to].receipts[credior_id];
            }
            
            //新增一张A欠C50的账单
            companies[before_addr].debts.push(
                Debt(companies[before_addr].debts.length, companies[to].receipts.length, to, amount)    
            );
            companies[to].receipts.push(
                Receipt(companies[to].receipts.length, companies[before_addr].debts.length, before_addr, amount)  
            );
            emit Transfer_bill(msg.sender, to, from_id, to_id, amount);
        }
    }
    
    //融资
    function finance(address addr, uint amount) public {
        require(msg.sender == bank.addr, "only bank could finance");
        //计算公司全部应收账款作为信用额度
        uint credit = 0;
        for(uint i = 0; i < companies[addr].receipts.length; ++i)
        {
            credit += companies[addr].receipts[i].amount;
        }
        require(credit >= amount, "credit of company is not enough");
        companies[addr].asset += amount;
        emit Finance(addr, amount);
    }
    
    //结算
    function settle(address to, uint bill_id) {
        uint credior_id = companies[msg.sender].debts[bill_id].credior_id;
        address creditor = companies[msg.sender].debts[bill_id].creditor;
        uint amount = companies[msg.sender].debts[bill_id].amount;
        require(companies[msg.sender].asset >= amount, "asset is not enough to pay");
        companies[msg.sender].asset -= amount;
        companies[creditor].asset += amount;
        delete companies[msg.sender].debts[bill_id];
        delete companies[to].receipts[credior_id];
        emit Settle(to, bill_id);
    }
}