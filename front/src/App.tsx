import './App.css'

import {Button, Link, Navbar, NavbarBrand, NavbarContent, NavbarItem} from "@nextui-org/react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import {faPlus} from "@fortawesome/free-solid-svg-icons";
import {useState} from "react";
import JobFormModal from "./components/jobFormModal.tsx";
import JobTable from "./components/jobTable.tsx";

function App() {
    const [ openFormJobModal, setOpenFormJobModal ] = useState(false);
    return (
        <>
            <Navbar isBordered={true}>
                <NavbarBrand>
                    <p className="font-bold text-inherit">Jobs Dashboard</p>
                </NavbarBrand>

                <NavbarContent justify="end">
                    <NavbarItem>
                        <Button
                            as={Link}
                            color="primary"
                            variant="flat"
                            onClick={() => setOpenFormJobModal(!openFormJobModal)}
                            startContent={<FontAwesomeIcon icon={faPlus} />}>
                            Add Job
                        </Button>
                    </NavbarItem>
                </NavbarContent>
            </Navbar>
            <JobFormModal isOpen={openFormJobModal} onOpenChange={setOpenFormJobModal} />
            <div className="p-4 container mx-auto max-w-screen-lg">
                <h1 className="text-3xl font-bold">Jobs</h1>
                <JobTable />
            </div>
        </>
    )
}

export default App
