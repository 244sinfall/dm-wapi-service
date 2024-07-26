import './styles.css'

type CheckTableInfoProps = {
    actualDate: string,
    checkCount: number
}

const CheckTableInfo = (props: CheckTableInfoProps) => {
    return (
        <div className="check-table-info">
            <p>Данные актуальны на: {props.actualDate}</p>
            <p>Чеков в БД: {props.checkCount}</p>
        </div>
    );
};

export default CheckTableInfo;